node {
  stage('SCM') {
    checkout scm
  }
  stage('SonarQube Analysis') {
    def scannerHome = tool 'SonarScanner';
    withSonarQubeEnv('mitraku-be-sonarqube') {
      sh "${scannerHome}/bin/sonar-scanner"
    }
  }
}