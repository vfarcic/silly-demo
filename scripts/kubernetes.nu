#!/usr/bin/env nu

def --env "main create kubernetes" [
    hyperscaler: string,
    --name = "dot",
    --min_nodes = 2,
    --max_nodes = 4,
    --node_size = "small" # Supported values: small, medium, large
    --auth = true
    --enable_ingress = true
] {

    $env.KUBECONFIG = $"($env.PWD)/kubeconfig-($name).yaml"
    $"export KUBECONFIG=($env.KUBECONFIG)\n" | save --append .env
    $"export KUBECONFIG_($name | str upcase)=($env.KUBECONFIG)\n" | save --append .env

    if $hyperscaler == "google" {

        if $auth {
            gcloud auth login
        }

        mut project_id = ""
        if PROJECT_ID in $env and not $auth {
            $project_id = $env.PROJECT_ID
        } else {
            $project_id = $"dot-(date now | format date "%Y%m%d%H%M%S")"
            $env.PROJECT_ID = $project_id
            $"export PROJECT_ID=($project_id)\n" | save --append .env

            gcloud projects create $project_id

            start $"https://console.cloud.google.com/marketplace/product/google/container.googleapis.com?project=($project_id)"
    
            print $"(ansi yellow_bold)
ENABLE(ansi reset) the API.
Press any key to continue.
"
            input
        }

        mut vm_size = "e2-standard-2"
        if $node_size == "medium" {
            $vm_size = "e2-standard-4"
        } else if $node_size == "large" {
            $vm_size = "e2-standard-8"
        }

        (
            gcloud container clusters create $name --project $project_id
                --zone us-east1-b --machine-type $vm_size
                --enable-autoscaling --num-nodes $min_nodes
                --min-nodes $min_nodes --max-nodes $max_nodes
                --enable-network-policy --no-enable-autoupgrade
        )

        (
            gcloud container clusters get-credentials $name
                --project $project_id --zone us-east1-b
        )

    } else if $hyperscaler == "aws" {

        (
            create eks  --name $name --node_size $node_size
                --min_nodes $min_nodes --max_nodes $max_nodes
        )

    } else if $hyperscaler == "azure" {

        mut tenant_id = ""
        let location = "eastus"

        if AZURE_TENANT in $env {
            $tenant_id = $env.AZURE_TENANT
        } else {
            $tenant_id = input $"(ansi green_bold)Enter Azure Tenant ID: (ansi reset)"
        }

        if $auth {
            az login --tenant $tenant_id
        }

        mut resource_group = ""
        if RESOURCE_GROUP in $env {
            $resource_group = $env.RESOURCE_GROUP
        } else {
            $resource_group = $"dot-(date now | format date "%Y%m%d%H%M%S")"
            $env.RESOURCE_GROUP = $resource_group
            $"export RESOURCE_GROUP=($resource_group)\n" | save --append .env
            az group create --name $resource_group --location $location
        }
        mut vm_size = "Standard_B2s"
        if $node_size == "medium" {
            $vm_size = "Standard_B4ms"
        } else if $node_size == "large" {
            $vm_size = "Standard_B8ms"
        }

        (
            az aks create --resource-group $resource_group --name $name
                --node-count $min_nodes --min-count $min_nodes
                --max-count $max_nodes
                --node-vm-size $vm_size
                --enable-managed-identity --generate-ssh-keys
                --enable-cluster-autoscaler --yes
        )

        (
            az aks get-credentials --resource-group $resource_group
                --name $name --file $env.KUBECONFIG
        )

    } else if $hyperscaler == "kind" {

        mut config = {
            kind: "Cluster"
            apiVersion: "kind.x-k8s.io/v1alpha4"
            name: $name
            nodes: [{
                role: "control-plane"
            }]
        }

        if $enable_ingress {
            $config = $config | merge {
                nodes: [{
                    role: "control-plane"
                    kubeadmConfigPatches: ['kind: InitConfiguration
nodeRegistration:
  kubeletExtraArgs:
    node-labels: "ingress-ready=true"'
                    ]
                    extraPortMappings: [{
                        containerPort: 80
                        hostPort: 80
                        protocol: "TCP"
                    }, {
                        containerPort: 443
                        hostPort: 443
                        protocol: "TCP"
                    }]
                }]
            }
        }
        
        $config | to yaml | save $"kind.yaml" --force

        kind create cluster --config kind.yaml
    
    } else {

        print $"(ansi red_bold)($hyperscaler)(ansi reset) is not a supported."
        exit 1

    }

    $env.KUBECONFIG

}

def "main destroy kubernetes" [
    hyperscaler: string
    --name = "dot"
    --delete_project = true
] {

    if $hyperscaler == "google" {

        rm --force kubeconfig.yaml

        (
            gcloud container clusters delete $name
                --project $env.PROJECT_ID --zone us-east1-b --quiet
        )

        if $delete_project {
            gcloud projects delete $env.PROJECT_ID --quiet
        }
    
    } else if $hyperscaler == "aws" {

        (
            eksctl delete addon --name aws-ebs-csi-driver
                --cluster $name --region us-east-1
        )

        (
            eksctl delete nodegroup --name primary
                --cluster $name --drain=false
                --region us-east-1 --parallel 10 --wait
        )

        (
            eksctl delete cluster
                --config-file $"eksctl-config-($name).yaml"
                --wait
        )

    } else if $hyperscaler == "azure" {

        (
            az aks delete --resource-group $env.RESOURCE_GROUP
                --name $name --yes
        )

        if $delete_project {

            az group delete --name $env.RESOURCE_GROUP --yes

        }

    } else if $hyperscaler == "kind" {

        kind delete cluster --name $name

    }

    rm --force kubeconfig.yaml

}

def "main create kubernetes_creds" [
    --source_kuberconfig = "kubeconfig.yaml"
    --destination_kuberconfig = "kubeconfig_new.yaml"
] {

    {
        apiVersion: "v1"
        kind: "ServiceAccount"
        metadata: {
            name: "creds"
            namespace: "kube-system"
        }
    } | to yaml | kubectl --kubeconfig $source_kuberconfig apply --filename -

    {
        apiVersion: "v1"
        kind: "Secret"
        metadata: {
            name: "creds"
            namespace: "kube-system"
            annotations: {
                "kubernetes.io/service-account.name": "creds"
            }
        }
        type: "kubernetes.io/service-account-token"
    } | to yaml | kubectl --kubeconfig $source_kuberconfig apply --filename -

    {
        apiVersion: "rbac.authorization.k8s.io/v1"
        kind: "ClusterRoleBinding"
        metadata: {
            name: "creds"
        }
        subjects: [{
            kind: "ServiceAccount"
            name: "creds"
            namespace: "kube-system"
        }]
        roleRef: {
            kind: "ClusterRole"
            name: "cluster-admin"
            apiGroup: "rbac.authorization.k8s.io"
        }
    }
        | to yaml
        | kubectl --kubeconfig $source_kuberconfig apply --filename -

    let kube_ca_data = open $source_kuberconfig
        | get clusters.0.cluster.certificate-authority-data

    let kube_url = open $source_kuberconfig
        | get clusters.0.cluster.server

    let token_encoded = (
        kubectl
            --kubeconfig $source_kuberconfig
            --namespace kube-system
            get secret creds --output yaml
    )
        | from yaml
        | get data.token

    let token = ($token_encoded | decode base64 | decode)

    {
        apiVersion: "v1"
        kind: "Config"
        clusters: [{
            name: "default-cluster"
            cluster: {
                certificate-authority-data: $kube_ca_data
                server: $"($kube_url):443"
            }
        }]
        contexts: [{
            name: "default-context"
            context: {
                cluster: "default-cluster"
                namespace: "default"
                user: "default-user"
            }
        }]
        current-context: "default-context"
        users: [{
            name: "default-user"
            user: {
                token: $token
            }
        }]
    } | to yaml | save $source_kuberconfig --force

}

def "create eks" [
    --name = "dot",
    --min_nodes = 2,
    --max_nodes = 4,
    --node_size = "small" # Supported values: small, medium, large
] {

    mut aws_access_key_id = ""
    if AWS_ACCESS_KEY_ID in $env {
        $aws_access_key_id = $env.AWS_ACCESS_KEY_ID
    } else {
        $aws_access_key_id = input $"(ansi green_bold)Enter AWS Access Key ID: (ansi reset)"
    }
    $"export AWS_ACCESS_KEY_ID=($aws_access_key_id)\n"
        | save --append .env

    mut aws_secret_access_key = ""
    if AWS_SECRET_ACCESS_KEY in $env {
        $aws_secret_access_key = $env.AWS_SECRET_ACCESS_KEY
    } else {
        $aws_secret_access_key = input $"(ansi green_bold)Enter AWS Secret Access Key: (ansi reset)" --suppress-output
    }
    $"export AWS_SECRET_ACCESS_KEY=($aws_secret_access_key)\n"
        | save --append .env

    mut aws_account_id = ""
    if AWS_ACCOUNT_ID in $env {
        $aws_account_id = $env.AWS_ACCOUNT_ID
    } else {
        $aws_account_id = input $"(ansi green_bold)Enter AWS Account ID: (ansi reset)"
    }
    $"export AWS_ACCOUNT_ID=($aws_account_id)\n"
        | save --append .env

    $"[default]
aws_access_key_id = ($aws_access_key_id)
aws_secret_access_key = ($aws_secret_access_key)
" | save aws-creds.conf --force

    mut vm_size = "t3.medium"
    if $node_size == "medium" {
        $vm_size = "t3.xlarge"
    } else if $node_size == "large" {
        $vm_size = "t3.2xlarge"
    }

    {
        apiVersion: "eksctl.io/v1alpha5"
        kind: "ClusterConfig"
        metadata: {
            name: $name
            region: "us-east-1"
            version: "1.31"
        }
        managedNodeGroups: [{
            name: "primary"
            instanceType: $vm_size
            minSize: $min_nodes
            maxSize: $max_nodes
            iam: {
                withAddonPolicies: {
                    autoScaler: true
                    ebs: true
                }
            }
        }]
    } | to yaml | save $"eksctl-config-($name).yaml" --force

    (
        eksctl create cluster
            --config-file $"eksctl-config-($name).yaml"
            --kubeconfig $env.KUBECONFIG
    )

    (
        eksctl create addon --name aws-ebs-csi-driver
            --cluster $name
            --service-account-role-arn $"arn:aws:iam::($aws_account_id):role/AmazonEKS_EBS_CSI_DriverRole"
            --region us-east-1 --force
    )

    (
        kubectl patch storageclass gp2
            --patch '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
    )

}