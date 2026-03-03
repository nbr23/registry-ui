pipeline {
    agent any

    options {
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout'){
            steps {
                checkout scm
            }
        }
        stage('Prep buildx') {
            steps {
                script {
                    env.BUILDX_BUILDER = getBuildxBuilder();
                }
            }
        }
        stage('Build full Image') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKERHUB_CREDENTIALS_USR', passwordVariable: 'DOCKERHUB_CREDENTIALS_PSW')]) {
                    sh 'docker login -u $DOCKERHUB_CREDENTIALS_USR -p "$DOCKERHUB_CREDENTIALS_PSW"'
                }
                sh """
                    docker buildx build \
                        --pull \
                        --builder \$BUILDX_BUILDER  \
                        --platform linux/arm64,linux/amd64 \
                        -t nbr23/registry-ui:latest \
                        -t nbr23/registry-ui:`git rev-parse --short HEAD` \
                        ${ "$GIT_BRANCH" == "master" ? "--push" : ""} .
                    """
            }
        }
        stage('Sync github repo') {
            when { branch 'master' }
            steps {
                ghSync()
            }
        }
    }
}
