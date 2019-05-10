# Payment GoAPI

Payment GoAPI is deployed on internet-facing Load Balancer with 2 docker private instances connected to the private mongodb cluster, nodes of which are on different VPCs.
* Application instances are on Virginia region and is accessable through Load Balancer public DNS: Payment-API-1319694113.us-east-1.elb.amazonaws.com
* Mongodb cluster has 3 nodes on N.California and 2 nodes on Oregon. They're not publicly accessable