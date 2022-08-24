.PHONY: package
package:
	helm package ./charts/mpijob -d ./charts/mpijob
