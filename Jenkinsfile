podTemplate(
    label: 'kubernetes',
    containers: [
        containerTemplate(name: 'maven', image: 'maven:alpine', ttyEnabled: true, command: 'cat'),
        containerTemplate(name: 'golang', image: 'golang:alpine', ttyEnabled: true, command: 'cat')
    ]
) {
    node('kubernetes') {
        container('maven') {
            stage('build') {
                sh "sleep 5"
                sh 'mvn --version'
            }
            stage('unit-test') {
                sh "sleep 5"
                sh 'java -version'
            }
        }
        container('golang') {
            stage('deploy') {
                sh "sleep 5"
                sh 'go version'
            }
        }
    }
}
