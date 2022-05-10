pipeline {
    environment {
      DOCKER = credentials('docker_hub')
    }
    agent any
        stages {
            stage('Build') {
                parallel {
                    stage('Express Image') {
                        steps {
                            sh 'docker build -f Dockerfile \
                            -t omollo/ktechnics-api-prod:latest .'
                        }
                    }                    
                }
                post {
                    failure {
                        echo 'This build has failed. See logs for details.'
                    }
                }
            }
            stage('Test') {
                steps {
                    echo 'This is the Testing Stage'
                }
            }
            stage('DEPLOY') {
                when {
                    branch 'main'  //only run these steps on the main branch
                }
                steps {
                    // sh 'docker tag ktechnics/ktechnics-api-dev:latest omollo/ktechnics/ktechnics-api-prod:latest'
                    sh 'docker login -u "omollo" -p "safcom2012" docker.io'
                    sh 'docker push omollo/ktechnics-api-prod:latest'
                }
            }

            stage('PUBLISH') {
                when {
                    branch 'main'  //only run these steps on the main branch
                }
                steps {
                    // sh 'docker swarm leave -f'
                    // sh 'docker run -d -p 8081:8081 --rm --name ekas-portal ktechnics/ktechnics-api-dev'
                    sh 'docker swarm init --advertise-addr 199.192.18.137'
                    sh 'docker stack deploy -c docker-compose.yml ktechnics-api-prod'
                }

            }

            // stage('REPORTS') {
            //     steps {
            //         junit 'reports.xml'
            //         archiveArtifacts(artifacts: 'reports.xml', allowEmptyArchive: true)
            //         // archiveArtifacts(artifacts: 'ktechnics/ktechnics-api-prod-golden.tar.gz', allowEmptyArchive: true)
            //     }
            // }

            stage('CLEAN-UP') {
                steps {
                    // sh 'docker stop ktechnics/ktechnics-api-dev'
                    sh 'docker system prune -f'
                    deleteDir()
                }
            }
        }
    }