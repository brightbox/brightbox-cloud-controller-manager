#!/usr/bin/env ruby

$LOAD_PATH.unshift File.expand_path('kubernetes-client/lib', __dir__)

require 'kubernetes'
require 'active_support/core_ext/hash/keys'
require 'active_support/core_ext/string/inflections'
require 'yaml'

first_release = 14
last_release = 17
config = {
  prefix: 'cloud-controller-build',
  image: 'gcr.io/kaniko-project/executor:latest',
  requestsCpu: '1700m',
  requestsMemory: '1Gi',
  limitsCpu: '2',
  limitsMemory: '1890Mi',
  gitRepo: 'git://github.com/brightbox/brightbox-cloud-controller-manager.git',
  dockerTarget: 'brightbox/brightbox-cloud-controller-manager',
  secretName: 'regcred',
  holdTime: 600
}

def make(target, attribute, attrs = {})
  create_object = lambda { |class_name, attrs|
    "Kubernetes::#{class_name}".constantize.new(attrs)
  }
  new_class = target.class.openapi_types[attribute.to_sym]
  target.public_send(
    :"#{attribute}=",
    case new_class
    when nil
      raise ArgumentError, attribute
    when /^Array<(?<class_name>\S+)>$/
      [create_object.call(Regexp.last_match[:class_name], attrs)]
    when /^Hash</
      attrs
    else
      create_object.call(new_class, attrs)
    end
  )
end

def create_job_manifest(release, config)
  version = `git describe --always release-#{release}`.strip.sub(/^v([0-9]+\.[0-9]+\.[0-9]+).*$/, '\1')
  name = "#{config[:prefix]}-#{version.tr('.', '-')}"
  job = Kubernetes::IoK8sApiBatchV1Job.new(
    api_version: 'batch/v1',
    kind: 'Job'
  )
  make(job, :metadata,
       name: name,
       labels: { build: config[:prefix] })
  make(job, :spec,
       ttl_seconds_after_finished: config[:holdTime])
  make(job.spec, :template)
  make(job.spec.template, :spec,
       restart_policy: 'Never')
  spec = job.spec.template.spec
  make(spec, :volumes,
       name: config[:secretName])
  volume = spec.volumes.first
  make(volume, :secret,
       secret_name: config[:secretName])
  make(volume.secret, :items,
       key: '.dockerconfigjson',
       path: '.docker/config.json')
  make(spec, :containers,
       name: name,
       image: config[:image],
       args: [
         '--dockerfile=Dockerfile',
         "--context=#{config[:gitRepo]}#refs/heads/release-#{release}",
         "--destination=#{config[:dockerTarget]}:#{version}"
       ])
  container = spec.containers.first
  make(container, :resources,
       requests: {
         memory: config[:requestsMemory],
         cpu: config[:requestsCpu]
       },
       limits: {
         memory: config[:limitsMemory],
         cpu: config[:limitsCpu]
       })
  make(container, :volume_mounts,
       name: config[:secretName],
       mount_path: '/root')
  job.to_hash.deep_stringify_keys.to_yaml
end

jobs = (first_release..last_release).inject('') do |memo, release|
  memo << create_job_manifest("1.#{release}", config)
end

puts jobs
# stdout_str, status = Open3.capture2('kubectl apply -f -', stdin_data: jobs)
# puts stdout_str
# exit status.exitstatus
