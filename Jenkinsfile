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
        
        // Set GOPATH for the build
        GOPATH = "${WORKSPACE}\\go"
        PATH = "${env.PATH};${GOPATH}\\bin"
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
                    // Display Go version
                    bat 'go version'
                    
                    // Create GOPATH directory
                    bat 'if not exist %GOPATH%\\bin mkdir %GOPATH%\\bin'
                    
                    // Install gotestsum for JUnit reports
                    bat 'go install gotest.tools/gotestsum@latest'
                    
                    // Verify gotestsum installed
                    bat 'dir %GOPATH%\\bin\\gotestsum.exe'
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
                    
                    # Run tests and generate coverage
                    go test -v -coverprofile=coverage.out -covermode=atomic ./...
                    
                    # Generate JUnit XML report using gotestsum
                    %GOPATH%\\bin\\gotestsum --junitfile test-report.xml -- -v -coverprofile=coverage.out ./...
                    
                    echo "Generating HTML coverage report..."
                    go tool cover -html=coverage.out -o coverage.html
                    
                    echo "Tests completed successfully"
                '''
            }
            post {
                always {
                    // Publish test results if file exists
                    script {
                        if (fileExists('test-report.xml')) {
                            junit 'test-report.xml'
                            echo "✅ Test results published"
                        } else {
                            echo "⚠️ No test-report.xml found"
                        }
                    }
                    
                    // Archive coverage report
                    archiveArtifacts artifacts: 'coverage.html', fingerprint: true
                }
            }
        }
        
        stage('SonarQube Analysis') {
            steps {
                script {
                    // Get SonarQube scanner tool with correct name
                    def scannerHome = tool 'SonarQubeScanner'
                    
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