pipeline {
    agent any

    environment {
        DOCKER_USER = "sudssy"
        IMAGE_NAME  = "${DOCKER_USER}/blog-backend"
        IMAGE_TAG   = "v${BUILD_NUMBER}"
        KUBECTL     = "/usr/bin/kubectl"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/sushant2122/assignment'
            }
        }

        stage('Build image') {
            steps {
                sh "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} ."
            }
        }

        stage('Push to DockerHub') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-creds',
                    usernameVariable: 'DUSER',
                    passwordVariable: 'DPASS'
                )]) {
                    sh "echo $DPASS | docker login -u $DUSER --password-stdin"
                    sh "docker push ${IMAGE_NAME}:${IMAGE_TAG}"
                }
            }
        }

        stage('Update Kubernetes') {
            steps {
                sh "${KUBECTL} rollout restart deployment/postgres -n blog"
                sh "${KUBECTL} rollout status deployment/postgres -n blog --timeout=60s"
                sh "${KUBECTL} set image deployment/backend backend=${IMAGE_NAME}:${IMAGE_TAG} -n blog"
                sh "${KUBECTL} rollout status deployment/backend -n blog --timeout=120s"
            }
        }
    }

    post {
        always { sh "docker logout" }
    }
}
