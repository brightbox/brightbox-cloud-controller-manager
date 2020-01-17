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

  jobs: {
    local release(version) = std.join('.', std.split(version, '.')[:2]),
    local releaseName(version) = "%s-%s" % [$._config.jobs.name, std.strReplace(version, ".", "-")],
     [std.toString(x)]: $.batch.v1.job.new() +
	{
	  "metadata": {
	    "name": releaseName(x),
	    "labels": {
	      "build": $._config.jobs.name
	    }
	  },
	  "spec": {
	    "ttlSecondsAfterFinished": $._config.jobs.holdTime,
	    "template": {
	      "spec": {
		"containers": [
		  $.core.v1.container.new(releaseName(x), $._config.jobs.image) +
		  $.util.resourcesRequests($._config.jobs.requestsCpu, $._config.jobs.requestsMemory) +
		  $.util.resourcesLimits($._config.jobs.limitsCpu, $._config.jobs.limitsMemory) +
		  {
		    "args": [
		      "--dockerfile=Dockerfile",
		      "--context=" + $._config.jobs.gitRepo + "#refs/heads/release-" + release(x),
		      "--destination=" + $._config.jobs.dockerTarget + ":" + x,
		    ],
		    "volumeMounts": [
		      $.core.v1.volumeMount.new($._config.jobs.secretName, "/root")
		    ]
		  }
		],
		"restartPolicy": "Never",
		"volumes": [
		  {
		    "name": $._config.jobs.secretName,
		    "secret": {
		      "defaultMode": 256,
		      "secretName": $._config.jobs.secretName,
		      "items": [
			{
			  "key": ".dockerconfigjson",
			  "path": ".docker/config.json"
			}
		      ]
		    }
		  }
		]
	      }
	    }
	  }
	}
     for x in versions
  }
}
