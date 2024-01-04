deploy:
 	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))
 	$(eval context=$(or $(context),k0s-dev-cluster))

	helm --kubeconfig ~/.kube --kube-context $(context) upgrade \
		--install vault \
		--values=./deploy/helm/vault/values.yaml \
		--values=./deploy/helm/vault/values_$(env).yaml \
		./deploy/helm/vault

deploy_vault:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval context=$(or $(context),k0s-dev-cluster))
	$(eval platform=$(or $(platform),linux/amd64))

	helm --kube-context $(context) upgrade \
		--install vault \
		--values=./deploy/helm/vault/values.yaml \
		--values=./deploy/helm/vault/values_$(env).yaml \
		./deploy/helm/vault

destroy_vault:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval context=$(or $(context),k0s-dev-cluster))
	$(eval platform=$(or $(platform),linux/amd64))

	helm --kube-context $(context) uninstall vault

.PHONY: deploy_vault