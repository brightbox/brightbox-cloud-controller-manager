#!/usr/bin/env ruby

$LOAD_PATH.unshift File.expand_path('kubernetes-client/lib', __dir__) 

require 'kubernetes'
require 'active_support/core_ext/hash/keys'
require 'open3'
require 'yaml'

first_release=14
last_release=17
config = {
  prefix:"cloud-controller-build",
  image:"gcr.io/kaniko-project/executor:latest",
  requestsCpu:"1700m",
  requestsMemory:"1Gi",
  limitsCpu:"2",
  limitsMemory:"1890Mi",
  gitRepo:"git://github.com/brightbox/brightbox-cloud-controller-manager.git",
  dockerTarget:"brightbox/brightbox-cloud-controller-manager",
  secretName:"regcred",
  holdTime:600,
}

def create_job_manifest(release, config)
  version = `git describe --always release-#{release}`.strip.sub(/^v([0-9]+\.[0-9]+\.[0-9]+).*$/, '\1')
  name = "#{config[:prefix]}-#{version.gsub('.', '-')}"
  volume_mount = Kubernetes::IoK8sApiCoreV1VolumeMount.new(
    name: config[:secretName], mount_path: "/root"
  )
  volume = Kubernetes::IoK8sApiCoreV1Volume.new(
    name: config[:secretName],
    secret: Kubernetes::IoK8sApiCoreV1SecretVolumeSource.new(
      secret_name: config[:secretName],
      items: [
        Kubernetes::IoK8sApiCoreV1KeyToPath.new(
          key: ".dockerconfigjson", path: ".docker/config.json"
        )
      ]
    )
  )
  container = Kubernetes::IoK8sApiCoreV1Container.new(
    name: name,
    image: config[:image],
    resources: Kubernetes::IoK8sApiCoreV1ResourceRequirements.new(
      requests: {
        memory: config[:requestsMemory],
        cpu: config[:requestsCpu],
      },
      limits: {
        memory: config[:limitsMemory],
        cpu: config[:limitsCpu],
      },
    ),
    args: [
      "--dockerfile=Dockerfile",
      "--context=#{config[:gitRepo]}#refs/heads/release-#{release}",
      "--destination=#{config[:dockerTarget]}:#{version}"
    ],
    volume_mounts: [volume_mount]
  )
  job = Kubernetes::IoK8sApiBatchV1Job.new(
    api_version: "batch/v1",
    kind: "Job",
    metadata: Kubernetes::IoK8sApimachineryPkgApisMetaV1ObjectMeta.new(
      name: name,
      labels: { build: config[:prefix] }
    ),
    spec: Kubernetes::IoK8sApiBatchV1JobSpec.new(
      ttl_seconds_after_finished: config[:holdTime],
      template: Kubernetes::IoK8sApiCoreV1PodTemplateSpec.new(
        spec: Kubernetes::IoK8sApiCoreV1PodSpec.new(
          restart_policy: "Never",
          containers: [container],
          volumes: [volume],
        )
      )
    ),
  )
  job.to_hash.deep_stringify_keys.to_yaml
end

jobs=(first_release..last_release).inject("") do |memo, release|
  memo << create_job_manifest("1.#{release}", config)
end

puts "Running jobs on k8s"
puts jobs
#stdout_str, status = Open3.capture2('kubectl apply -f -', stdin_data: jobs)
#puts stdout_str
#exit status.exitstatus
