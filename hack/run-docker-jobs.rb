#!/usr/bin/env ruby

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
  {
    apiVersion: "batch/v1",
    kind: "Job",
    metadata: {
      name: name,
      labels: {
        build: config[:prefix]
      },
    },
    spec: {
      ttlSecondsAfterFinished: config[:holdTime],
      template: {
        spec: {
          containers: [
            {
              name: name,
              image: config[:image],
              resources: {
                requests: {
                  memory: config[:requestsMemory],
                  cpu: config[:requestsCpu],
                },
                limits: {
                  memory: config[:limitsMemory],
                  cpu: config[:limitsCpu],
                },
              },
              args: [
                "--dockerfile=Dockerfile",
                "--context=#{config[:gitRepo]}#refs/heads/release-#{release}",
                "--destination=#{config[:dockerTarget]}:#{version}"
              ],
              volumeMounts: [
                {
                  name: config[:secretName],
                  mountPath: "/root",
                },
              ],
            },
          ],
          restartPolicy: "Never",
          volumes: [
            {
              name: config[:secretName],
              secret: {
                secretName: config[:secretName],
                items: [
                  {
                    key: ".dockerconfigjson",
                    path: ".docker/config.json"
                  },
                ],
              },
            },
          ],
        },
      },
    },
  }.deep_transform_keys(&:to_s).to_yaml
end

jobs=(first_release..last_release).inject("") do |memo, release|
  memo << create_job_manifest("1.#{release}", config)
end

puts "Running jobs on k8s"
stdout_str, status = Open3.capture2('kubectl apply -f -', stdin_data: jobs)
puts stdout_str
exit status.exitstatus
