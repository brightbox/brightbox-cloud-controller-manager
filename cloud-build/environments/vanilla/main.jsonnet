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
     [std.toString(x)]:
	{
	  "apiVersion": "batch/v1",
	  "kind": "Job",
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
		  {
		    "name": releaseName(x),
		    "image": $._config.jobs.image,
		    "imagePullPolicy": "IfNotPresent",
		    "resources": {
		      "requests": {
			"memory": $._config.jobs.requestsMemory,
			"cpu": $._config.jobs.requestsCpu,
		      },
		      "limits": {
			"memory": $._config.jobs.limitsMemory,
			"cpu": $._config.jobs.limitsCpu,
		      }
		    },
		    "args": [
		      "--dockerfile=Dockerfile",
		      "--context=" + $._config.jobs.gitRepo + "#refs/heads/release-" + release(x),
		      "--destination=" + $._config.jobs.dockerTarget + ":" + x,
		    ],
		    "volumeMounts": [
		      {
			"name": $._config.jobs.secretName,
			"mountPath": "/root"
		      }
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
