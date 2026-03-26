pipeline {
    agent any
    
    tools {
        go 'go-1.25.1'  // Configure this in Jenkins Global Tool Configuration
    }
    
    environment {
        // SonarQube configuration
        SONAR_HOST_URL = 'http://localhost:9000'
        SONAR_TOKEN = credentials('sonarqube-token')
        
        // Go environment
        GO111MODULE = 'on'
        GOPATH = "${WORKSPACE}/go"
    }
    
    stages {
        stage('Checkout') {
            steps {
                // Replace with your actual repository
                git url: 'https://github.com/yourusername/go-crud-app.git', branch: 'main'
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
                    # Run tests with coverage
                    go test -v -coverprofile=coverage.out -covermode=atomic ./...
                    
                    # Generate JSON test report
                    go test -json ./... > report.json
                    
                    # Generate coverage HTML for debugging
                    go tool cover -html=coverage.out -o coverage.html
                '''
            }
            post {
                always {
                    // Publish test results
                    junit 'report.json'
                    
                    // Archive coverage report
                    archiveArtifacts artifacts: 'coverage.html', fingerprint: true
                }
            }
        }
        
        stage('Static Analysis') {
            steps {
                sh '''
                    # Run go vet
                    go vet ./...
                    
                    # Run staticcheck if installed
                    # staticcheck ./...
                    
                    # Run golint if installed
                    # golint ./...
                '''
            }
        }
        
        stage('SonarQube Analysis') {
            steps {
                withSonarQubeEnv('SonarQube') {
                    sh '''
                        sonar-scanner \
                          -Dsonar.projectKey=go-crud-app \
                          -Dsonar.projectName="Go CRUD Application" \
                          -Dsonar.sources=. \
                          -Dsonar.exclusions="**/*_test.go,**/vendor/*" \
                          -Dsonar.tests=. \
                          -Dsonar.test.inclusions="**/*_test.go" \
                          -Dsonar.go.coverage.reportPaths=coverage.out \
                          -Dsonar.go.tests.reportPaths=report.json
                    '''
                }
            }
        }
        
        stage('Quality Gate') {
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                    # Build the application
                    go build -o bin/go-crud-app .
                '''
            }
            post {
                success {
                    archiveArtifacts artifacts: 'bin/go-crud-app', fingerprint: true
                }
            }
        }
    }
    
    post {
        always {
            // Clean up
            cleanWs()
        }
        success {
            echo '✅ Pipeline passed! Check SonarQube dashboard for detailed analysis'
            echo "View results: ${SONAR_HOST_URL}/dashboard?id=go-crud-app"
        }
        failure {
            echo '❌ Pipeline failed. Check SonarQube quality gate results.'
        }
    }
}