pipeline {
    agent any

    environment {
        SONAR_SCANNER_HOME = tool 'SonarScanner'
    }

    stages {
        stage('Build and SonarQube Analysis') {
            steps {
                script {
                    // Clone the repository
                    checkout([$class: 'GitSCM', branches: [[name: '*/development']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[url: 'https://github.com/Mitra-Apps/be-user-service']]])

                    // Set up your Go environment (install dependencies, etc.)
                    sh 'go get -d -v ./...'
                    sh 'go install -v ./...'

                    // Build the Go code
                    sh 'go build -o my_go_app'

                    // Run SonarQube Scanner
                    withSonarQubeEnv('mitraku-be-sonarqube') {
                        sh "${SONAR_SCANNER_HOME}/bin/sonar-scanner"
                    }
                }
            }
        }
    }
}
