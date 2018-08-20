**Lil'Dns-Proxy**

This app is made in Golang using the following libraries:

	crypto/tls
	crypto/x509
	fmt
	log
	net

*Problem Statement*

The applications need a proxy to talk TLS with some DNS Resolvers like Cloudflare and Google to increase security.

*Approach*

This application was made in Go to be simple to run as a binary inside a previous layer of the application container.

If the container infrastructure already run in Kubernetes or Docker Swarm, the app is also provided with a Dockerfile to be deployed in the pod/swarm group.

*Security Concerns*

In a microservice architecture I personally recommend run this app in the same container as the application, to make all DNS queries leave the container already encrypted.

Running this in the same environment exposing the port 53 could led an attacker to be able to sniff the DNS requests.

*Improvements for the code.*

 - [ ] Remove hard-coded variables and make ENV_VARS
 - [ ] Remove the certificate from code and add request via Let's Encrypt or other provider
 - [ ] Add a UDP Listener
 - [ ] Implement best approach for the dial.tls to make it more clean and easy to maintain
 - [ ] Add a DEBUG option to be used via ENV_VARS.
