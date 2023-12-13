pipeline {
    agent any

    environment {
        // Define environment variables
        GO_VERSION = '1.21.3'
        DOCKER_COMPOSE_VERSION = '2.21.0'
        DOCKER_COMPOSE_FILE = 'docker-compose.yaml'
    }

    stages {
        stage('Checkout') {
            steps {
                // This stage checks out the source code from your version control system
                checkout scm
            }
        }

        stage('Build') {
            steps {
                // Use a Golang image to build and test your Go application
                container('golang:' + GO_VERSION) {
                    sh 'go mod download'
                    sh 'go mod tidy'
                    sh 'go mod vendor'
                    sh 'go build'
                    echo "INFO: Code has been built"
                }
            }
        }

        stage('Run Docker Compose') {
            steps {
                // Run Docker Compose to start your application and any required services
                script {
                    sh "docker rm -f be-user-service-user-service-1 || true"
                    def dockerComposeCmd = "docker compose up -d"
                    sh dockerComposeCmd
                    echo "INFO: Deployed"
                }
            }
        }
    }

    post {
        success {
            // This block is executed if the pipeline is successful
            echo 'Pipeline succeeded! Send notifications or perform additional tasks here.'
        }

        failure {
            // This block is executed if the pipeline fails
            echo 'Pipeline failed! Send notifications or perform additional tasks here.'
        }
    }
}