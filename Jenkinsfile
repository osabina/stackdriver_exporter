node ("docker") {
    checkoutRepo()

    stage ( "Build" ) {
        try {
            startGHCheck("Docker Build",currentBuild.absoluteUrl,"Building docker image...")
            withDocker {
                sh "scripts/dockerbuild.sh build"
            }
        passGHCheck("Docker Build",currentBuild.absoluteUrl,"Docker image was built successfully!")
      } catch (e) {
        echo "Exception thrown:\n ${e}"
        failGHCheck("Docker Build",currentBuild.absoluteUrl,"Docker image failed to build")
        throw e
      }
    }
}

if ( isBranch("main") ) {
    stage( "Deploy" ) {
        withDocker {
                sh "scripts/dockerbuild.sh push"
            }
        }
}