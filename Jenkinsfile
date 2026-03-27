pipeline {
    agent any
    
    environment {
        // SonarQube configuration
        SONAR_HOST_URL = 'http://host.docker.internal:9000'  // For Docker on Windows
        
        // Your SonarQube token from Jenkins credentials
        SONAR_TOKEN = credentials('sonarqube-global-token')
        
        // Go environment
        GO111MODULE = 'on'
        
        // Project configuration
        PROJECT_KEY = 'test-sonarqube'
        PROJECT_NAME = 'Test SonarQube Go App'
        
        // Use workspace-specific cache to avoid file lock issues
        GOCACHE = "${WORKSPACE}\\gocache"
        GOPATH = "${WORKSPACE}\\go"
        PATH = "${env.PATH};${GOPATH}\\bin"
    }
    
    stages {
        stage('SCM Checkout') {
            steps {
                // Clean workspace before checkout to prevent file lock issues
                cleanWs()
                checkout scm
                echo "✅ Code checked out to: ${WORKSPACE}"
                
                // Verify project structure
                bat 'echo "Project files:" && dir'
                bat 'echo "Test folder:" && if exist test dir test'
            }
        }
        
        stage('Setup Go Environment') {
            steps {
                script {
                    // Display Go version
                    bat 'go version'
                    
                    // Create necessary directories
                    bat 'if not exist %GOCACHE% mkdir %GOCACHE%'
                    bat 'if not exist %GOPATH%\\bin mkdir %GOPATH%\\bin'
                    
                    // Install gotestsum (optional, for JUnit reports)
                    bat 'go install gotest.tools/gotestsum@latest'
                    echo "✅ Go environment configured"
                }
            }
        }
        
        stage('Download Dependencies') {
            steps {
                bat '''
                    echo "Downloading dependencies..."
                    go mod download
                    go mod verify
                    echo "✅ Dependencies downloaded"
                '''
            }
        }
        
        stage('Run Tests with Coverage') {
            steps {
                bat '''
                    echo "Running tests with coverage..."
                    
                    # Run all tests and generate coverage
                    go test -v -coverprofile=coverage.out -covermode=atomic ./...
                    
                    # Check if coverage.out was generated
                    if exist coverage.out (
                        echo "✅ coverage.out generated"
                        go tool cover -func=coverage.out | findstr total
                    ) else (
                        echo "❌ coverage.out NOT generated"
                        exit 1
                    )
                    
                    echo "Generating HTML coverage report..."
                    go tool cover -html=coverage.out -o coverage.html
                    
                    echo "✅ Tests completed successfully"
                '''
            }
            post {
                always {
                    script {
                        // Publish JUnit test results if available
                        if (fileExists('test-report.xml')) {
                            junit 'test-report.xml'
                            echo "✅ Test results published"
                        }
                        
                        // Archive coverage reports
                        if (fileExists('coverage.out')) {
                            archiveArtifacts artifacts: 'coverage.out', fingerprint: true
                        }
                        if (fileExists('coverage.html')) {
                            archiveArtifacts artifacts: 'coverage.html', fingerprint: true
                        }
                    }
                }
            }
        }
        
        stage('SonarQube Analysis') {
            steps {
                script {
                    // Ensure coverage.out exists before running SonarQube
                    if (!fileExists('coverage.out')) {
                        error "coverage.out not found! Cannot run SonarQube analysis."
                    }
                    
                    // IMPORTANT: Wrap Docker command in withSonarQubeEnv
                    withSonarQubeEnv('SonarQube') {
                        bat """
                            echo "Running SonarQube analysis with Docker..."
                            docker run --rm ^
                              -e SONAR_HOST_URL="${SONAR_HOST_URL}" ^
                              -e SONAR_TOKEN=${SONAR_TOKEN} ^
                              -v "${WORKSPACE}:/usr/src" ^
                              sonarsource/sonar-scanner-cli ^
                              "-Dsonar.projectKey=${PROJECT_KEY}" ^
                              "-Dsonar.projectName=${PROJECT_NAME}" ^
                              "-Dsonar.projectVersion=1.0.0" ^
                              "-Dsonar.sources=." ^
                              "-Dsonar.exclusions=**/*_test.go,**/vendor/*" ^
                              "-Dsonar.tests=./test" ^
                              "-Dsonar.test.inclusions=**/*_test.go" ^
                              "-Dsonar.go.coverage.reportPaths=coverage.out" ^
                              "-Dsonar.sourceEncoding=UTF-8"
                            echo "✅ SonarQube analysis completed"
                        """
                    }
                }
            }
        }
        
        stage('Wait for Quality Gate') {
            steps {
                timeout(time: 10, unit: 'MINUTES') {
                    // This will now work because we used withSonarQubeEnv above
                    waitForQualityGate abortPipeline: true
                }
                echo "✅ Quality gate passed"
            }
        }
        
        stage('Build Application') {
            steps {
                bat '''
                    echo "Building application..."
                    go build -o bin\\test-sonarqube.exe .
                    echo "✅ Build completed"
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
            // Clean up workspace to prevent file lock issues
            cleanWs()
            echo "🧹 Workspace cleaned"
        }
        success {
            echo '✅ Pipeline completed successfully!'
            echo "📊 View SonarQube results: ${SONAR_HOST_URL}/dashboard?id=${PROJECT_KEY}"
        }
        failure {
            echo '❌ Pipeline failed. Please check the logs above.'
            echo "🔍 Check SonarQube quality gate: ${SONAR_HOST_URL}/dashboard?id=${PROJECT_KEY}"
            currentBuild.result = 'FAILURE'
        }
    }
}