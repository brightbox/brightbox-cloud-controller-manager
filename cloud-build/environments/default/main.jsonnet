(import "ksonnet-util/kausal.libsonnet") +
{
  _config:: {
    jobs: {
      name: "cloud-controller-build",
      holdTime: 600,
      image: "gcr.io/kaniko-project/executor:latest",
      requestsCpu: "1700m",
      requestsMemory: "1Gi",
      limitsCpu: "2",
      limitsMemory: "1890Mi",
      gitRepo: "git://github.com/brightbox/brightbox-cloud-controller-manager.git",
      dockerTarget: "brightbox/brightbox-cloud-controller-manager",
      secretName: "regcred",
      },
  },
  local versions = [x for x in std.split(importstr "versions" ,"\n") if x != ""],

  util+:: {
    local release(version) = std.join('.', std.split(version, '.')[:2]),
    releaseName(version):: "%s-%s" % [$._config.jobs.name, std.strReplace(version, '.', '-')],
    versionArgs(version):: [
      "--dockerfile=Dockerfile",
      "--context=" + $._config.jobs.gitRepo + "#refs/heads/release-" + release(version),
      "--destination=" + $._config.jobs.dockerTarget + ":" + version,
    ],
    jobSecretVolumeMount(name, path, defaultMode=256, volumeMountMixin={})::
      local container = $.core.v1.container,
            job = $.batch.v1.job,
            deployment = $.extensions.v1beta1.deployment,
            volumeMount = $.core.v1.volumeMount,
            volume = $.core.v1.volume,
            addMount(c) = c + container.withVolumeMountsMixin(
        volumeMount.new(name, path) +
        volumeMountMixin,
      );

      job.mapContainers(addMount) +
      job.mixin.spec.template.spec.withVolumesMixin([
        volume.fromSecret(name, name) +
        volume.mixin.secret.withDefaultMode(defaultMode) +
	volume.mixin.secret.withItems({key: ".dockerconfigjson", path: ".docker/config.json"}),
	])
  },

  jobs: {
    local job = $.batch.v1.job,
    local container = $.core.v1.container,
     [std.toString(x)] : job.new() +
       job.mixin.metadata.withName($.util.releaseName(x)) +
       job.mixin.metadata.withLabels({build: $._config.jobs.name}) +
       job.mixin.spec.withTtlSecondsAfterFinished($._config.jobs.holdTime) +
       job.mixin.spec.template.spec.withRestartPolicy("Never") +
       job.mixin.spec.template.spec.withContainers([
         container.new($.util.releaseName(x), $._config.jobs.image) +
	 container.withArgs($.util.versionArgs(x)) +
	 $.util.resourcesRequests($._config.jobs.requestsCpu, $._config.jobs.requestsMemory) +
	 $.util.resourcesLimits($._config.jobs.limitsCpu, $._config.jobs.limitsMemory)
       ]) +
       $.util.jobSecretVolumeMount($._config.jobs.secretName, "/root")
     for x in versions
  }
}
