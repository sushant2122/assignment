pipeline {
    agent any

    environment {
        DOCKER_USER = "yourusername"
        IMAGE_NAME  = "${DOCKER_USER}/blog-backend"
        IMAGE_TAG   = "v${BUILD_NUMBER}"
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
                sh """
                    kubectl set image deployment/backend \
                      backend=${IMAGE_NAME}:${IMAGE_TAG} \
                      -n blog
                """
            }
        }
    }

    post {
        always { sh "docker logout" }
    }
}
