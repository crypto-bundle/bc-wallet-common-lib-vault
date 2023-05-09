deploy:
	helm --kubeconfig ~/.kube/config --kube-context docker-desktop upgrade \
		--install vault \
		--values=./deploy/helm/vault/values.yaml \
		--values=./deploy/helm/vault/values_$(env).yaml \
		./deploy/helm/vault

.PHONY: deploy