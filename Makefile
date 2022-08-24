.PHONY: package
package:
	helm package ./charts/mpijob -d ./charts/mpijob

.PHONY: template
template:
	helm template test-name ./charts/mpijob
