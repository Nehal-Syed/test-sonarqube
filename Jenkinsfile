// Complete Jenkinsfile for Go project with SonarQube analysis
pipeline {
    agent any
    
    environment {
        // SonarQube configuration
        SONAR_HOST_URL = 'http://localhost:9000'
        
        // Go environment
        GO111MODULE = 'on'
        
        // Project configuration
        PROJECT_KEY = 'test-sonarqube'
        PROJECT_NAME = 'Test SonarQube Go App'
    }
    
    stages {
        stage('SCM Checkout') {
            steps {
                checkout scm
                echo "✅ Code checked out successfully"
            }
        }
        
        stage('Setup Go Environment') {
            steps {
                script {
                    // Check if Go is available, if not install it
                    def goVersion = sh(script: 'go version', returnStatus: true)
                    if (goVersion != 0) {
                        echo "⚠️ Go not found. Installing Go..."
                        sh '''
                            # Download and install Go
                            wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
                            sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz
                            export PATH=$PATH:/usr/local/go/bin
                            echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
                        '''
                    }
                    sh 'go version'
                }
            }
        }
        
        stage('Download Dependencies') {
            steps {
                sh '''
                    go mod download
                    go mod verify
                '''
            }
        }
        
        stage('Run Tests with Coverage') {
            steps {
                sh '''
                    # Run tests and generate coverage report
                    go test -v -coverprofile=coverage.out -covermode=atomic ./...
                    
                    # Generate JSON test report for Jenkins
                    go test -json ./... > test-report.json
                    
                    # Generate HTML coverage report
                    go tool cover -html=coverage.out -o coverage.html
                '''
            }
            post {
                always {
                    // Publish test results
                    junit 'test-report.json'
                    
                    // Archive coverage report
                    archiveArtifacts artifacts: 'coverage.html', fingerprint: true
                }
            }
        }
        
        stage('SonarQube Analysis') {
            steps {
                script {
                    // Get SonarQube scanner tool
                    def scannerHome = tool 'SonarScanner'
                    
                    withSonarQubeEnv('SonarQube') {
                        sh """
                            ${scannerHome}/bin/sonar-scanner \
                              -Dsonar.projectKey=${PROJECT_KEY} \
                              -Dsonar.projectName="${PROJECT_NAME}" \
                              -Dsonar.projectVersion=1.0.0 \
                              -Dsonar.sources=. \
                              -Dsonar.exclusions="**/*_test.go,**/vendor/*" \
                              -Dsonar.tests=. \
                              -Dsonar.test.inclusions="**/*_test.go" \
                              -Dsonar.go.coverage.reportPaths=coverage.out \
                              -Dsonar.go.tests.reportPaths=test-report.json \
                              -Dsonar.sourceEncoding=UTF-8
                        """
                    }
                }
            }
        }
        
        stage('Wait for Quality Gate') {
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }
        
        stage('Build Application') {
            steps {
                sh '''
                    # Build the Go binary
                    go build -o bin/test-sonarqube .
                    
                    # Verify build
                    ls -la bin/
                '''
            }
            post {
                success {
                    archiveArtifacts artifacts: 'bin/test-sonarqube', fingerprint: true
                }
            }
        }
    }
    
    post {
        always {
            // Clean up workspace
            cleanWs()
        }
        success {
            echo '✅ Pipeline completed successfully!'
            echo "📊 View SonarQube results: ${SONAR_HOST_URL}/dashboard?id=${PROJECT_KEY}"
        }
        failure {
            echo '❌ Pipeline failed. Please check the logs above.'
            echo "🔍 Check SonarQube quality gate: ${SONAR_HOST_URL}/dashboard?id=${PROJECT_KEY}"
        }
    }
}