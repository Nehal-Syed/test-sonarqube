// Windows-compatible Jenkinsfile for Go project with SonarQube
pipeline {
    agent any
    
    environment {
        // SonarQube configuration
        SONAR_HOST_URL = 'http://localhost:9000'
        
        // Go environment for Windows
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
                    // Check if Go is available on Windows
                    def goCheck = bat(script: 'go version', returnStatus: true)
                    if (goCheck != 0) {
                        echo "❌ Go is not installed or not in PATH"
                        echo "Please install Go from: https://golang.org/dl/"
                        error "Go is required for this pipeline"
                    }
                    
                    // Display Go version
                    bat 'go version'
                }
            }
        }
        
        stage('Download Dependencies') {
            steps {
                bat '''
                    echo "Downloading dependencies..."
                    go mod download
                    go mod verify
                    echo "Dependencies downloaded successfully"
                '''
            }
        }
        
        stage('Run Tests with Coverage') {
            steps {
                bat '''
                    echo "Running tests with coverage..."
                    go test -v -coverprofile=coverage.out -covermode=atomic ./...
                    
                    echo "Generating test reports..."
                    go test -json ./... > test-report.json
                    
                    echo "Generating HTML coverage report..."
                    go tool cover -html=coverage.out -o coverage.html
                    
                    echo "Tests completed successfully"
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
                        bat """
                            echo "Running SonarQube analysis..."
                            ${scannerHome}\\bin\\sonar-scanner.bat ^
                              -Dsonar.projectKey=${PROJECT_KEY} ^
                              -Dsonar.projectName="${PROJECT_NAME}" ^
                              -Dsonar.projectVersion=1.0.0 ^
                              -Dsonar.sources=. ^
                              -Dsonar.exclusions="**/*_test.go,**/vendor/*" ^
                              -Dsonar.tests=. ^
                              -Dsonar.test.inclusions="**/*_test.go" ^
                              -Dsonar.go.coverage.reportPaths=coverage.out ^
                              -Dsonar.go.tests.reportPaths=test-report.json ^
                              -Dsonar.sourceEncoding=UTF-8
                            echo "SonarQube analysis completed"
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
                bat '''
                    echo "Building application..."
                    go build -o bin\\test-sonarqube.exe .
                    echo "Build completed successfully"
                    
                    echo "Build artifacts:"
                    dir bin\\
                '''
            }
            post {
                success {
                    archiveArtifacts artifacts: 'bin/test-sonarqube.exe', fingerprint: true
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