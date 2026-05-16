pipeline {
    agent any

    environment {
        // Armazene sua chave do Gemini nos Credentials do Jenkins como Secret Text
        GEMINI_API_KEY = credentials('gemini-api-key')
    }

    stages {
        stage('Sincronizar Repositório') {
            steps {
                checkout scm
                sh 'git fetch --prune --unshallow || git fetch --tags'
            }
        }

        stage('Instalar doc-ai') {
            steps {
                echo 'Instalando ou atualizando a CLI doc-ai em Go...'
                sh 'go install github.com/user/doc-ai@latest'
            }
        }

        stage('Gerar Documentação Incremental') {
            steps {
                echo 'Analisando o último commit e atualizando documentação...'
                sh 'doc-ai incremental .'
            }
        }

        stage('Publicar Alterações no Bitbucket') {
            steps {
                sh '''
                    # Configura identificação do Bot do Jenkins
                    git config user.email "jenkins-bot@suaempresa.com"
                    git config user.name "Jenkins CI Bot"

                    if [ -n "$(git status --porcelain docs/ .docai/)" ]; then
                        git add docs/ .docai/
                        git commit -m "docs: atualização automática de cobertura [ci skip]"
                        git push origin HEAD:${BITBUCKET_BRANCH}
                    else
                        echo "Nenhuma alteração de documentação necessária."
                    fi
                '''
            }
        }
    }

    post {
        failure {
            echo "Falha ao gerar a documentação automática."
        }
    }
}
