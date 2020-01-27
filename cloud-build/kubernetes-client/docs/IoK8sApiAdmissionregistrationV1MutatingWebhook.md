# Kubernetes::IoK8sApiAdmissionregistrationV1MutatingWebhook

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**admission_review_versions** | **Array&lt;String&gt;** | AdmissionReviewVersions is an ordered list of preferred &#x60;AdmissionReview&#x60; versions the Webhook expects. API server will try to use first version in the list which it supports. If none of the versions specified in this list supported by API server, validation will fail for this object. If a persisted webhook configuration specifies allowed versions and does not include any versions known to the API Server, calls to the webhook will fail and be subject to the failure policy. | 
**client_config** | [**IoK8sApiAdmissionregistrationV1WebhookClientConfig**](IoK8sApiAdmissionregistrationV1WebhookClientConfig.md) |  | 
**failure_policy** | **String** | FailurePolicy defines how unrecognized errors from the admission endpoint are handled - allowed values are Ignore or Fail. Defaults to Fail. | [optional] 
**match_policy** | **String** | matchPolicy defines how the \&quot;rules\&quot; list is used to match incoming requests. Allowed values are \&quot;Exact\&quot; or \&quot;Equivalent\&quot;.  - Exact: match a request only if it exactly matches a specified rule. For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1, but \&quot;rules\&quot; only included &#x60;apiGroups:[\&quot;apps\&quot;], apiVersions:[\&quot;v1\&quot;], resources: [\&quot;deployments\&quot;]&#x60;, a request to apps/v1beta1 or extensions/v1beta1 would not be sent to the webhook.  - Equivalent: match a request if modifies a resource listed in rules, even via another API group or version. For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1, and \&quot;rules\&quot; only included &#x60;apiGroups:[\&quot;apps\&quot;], apiVersions:[\&quot;v1\&quot;], resources: [\&quot;deployments\&quot;]&#x60;, a request to apps/v1beta1 or extensions/v1beta1 would be converted to apps/v1 and sent to the webhook.  Defaults to \&quot;Equivalent\&quot; | [optional] 
**name** | **String** | The name of the admission webhook. Name should be fully qualified, e.g., imagepolicy.kubernetes.io, where \&quot;imagepolicy\&quot; is the name of the webhook, and kubernetes.io is the name of the organization. Required. | 
**namespace_selector** | [**IoK8sApimachineryPkgApisMetaV1LabelSelector**](IoK8sApimachineryPkgApisMetaV1LabelSelector.md) |  | [optional] 
**object_selector** | [**IoK8sApimachineryPkgApisMetaV1LabelSelector**](IoK8sApimachineryPkgApisMetaV1LabelSelector.md) |  | [optional] 
**reinvocation_policy** | **String** | reinvocationPolicy indicates whether this webhook should be called multiple times as part of a single admission evaluation. Allowed values are \&quot;Never\&quot; and \&quot;IfNeeded\&quot;.  Never: the webhook will not be called more than once in a single admission evaluation.  IfNeeded: the webhook will be called at least one additional time as part of the admission evaluation if the object being admitted is modified by other admission plugins after the initial webhook call. Webhooks that specify this option *must* be idempotent, able to process objects they previously admitted. Note: * the number of additional invocations is not guaranteed to be exactly one. * if additional invocations result in further modifications to the object, webhooks are not guaranteed to be invoked again. * webhooks that use this option may be reordered to minimize the number of additional invocations. * to validate an object after all mutations are guaranteed complete, use a validating admission webhook instead.  Defaults to \&quot;Never\&quot;. | [optional] 
**rules** | [**Array&lt;IoK8sApiAdmissionregistrationV1RuleWithOperations&gt;**](IoK8sApiAdmissionregistrationV1RuleWithOperations.md) | Rules describes what operations on what resources/subresources the webhook cares about. The webhook cares about an operation if it matches _any_ Rule. However, in order to prevent ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks from putting the cluster in a state which cannot be recovered from without completely disabling the plugin, ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks are never called on admission requests for ValidatingWebhookConfiguration and MutatingWebhookConfiguration objects. | [optional] 
**side_effects** | **String** | SideEffects states whether this webhook has side effects. Acceptable values are: None, NoneOnDryRun (webhooks created via v1beta1 may also specify Some or Unknown). Webhooks with side effects MUST implement a reconciliation system, since a request may be rejected by a future step in the admission change and the side effects therefore need to be undone. Requests with the dryRun attribute will be auto-rejected if they match a webhook with sideEffects &#x3D;&#x3D; Unknown or Some. | 
**timeout_seconds** | **Integer** | TimeoutSeconds specifies the timeout for this webhook. After the timeout passes, the webhook call will be ignored or the API call will fail based on the failure policy. The timeout value must be between 1 and 30 seconds. Default to 10 seconds. | [optional] 

## Code Sample

```ruby
require 'Kubernetes'

instance = Kubernetes::IoK8sApiAdmissionregistrationV1MutatingWebhook.new(admission_review_versions: null,
                                 client_config: null,
                                 failure_policy: null,
                                 match_policy: null,
                                 name: null,
                                 namespace_selector: null,
                                 object_selector: null,
                                 reinvocation_policy: null,
                                 rules: null,
                                 side_effects: null,
                                 timeout_seconds: null)
```


