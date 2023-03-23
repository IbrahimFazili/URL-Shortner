import subprocess

hosts = ["10.11.1.107", "10.11.2.107", "10.11.3.107"]
swarms = []
manager = hosts[0]
join_swarm_command = ""

# CASSANDRA
def getHostId(nodeIP):
    byteInput = subprocess.Popen('docker exec -ti cassandra-node nodetool status'.split(" "), stdout=subprocess.PIPE)
    output = subprocess.check_output('grep {}'.format(nodeIP).split(" "), stdin=byteInput.stdout).decode('UTF-8').split(" ")
    byteInput.wait()
    hostId = output[-3]
    print(hostId)
    return hostId

def removeCassandra():
    subprocess.call(['cassandra/stopCluster.sh', *hosts])

def startCassandra():
    subprocess.call(['cassandra/startCluster.sh', *hosts])

def addCassandraNode(nodeIP):
    print("----------------\nAdding {} to Cassandra cluster\n-----------------".format(nodeIP))
    subprocess.call(['cassandra/create_volumes.sh', nodeIP])
    subprocess.call(['cassandra/addNodeToCluster.sh', manager, nodeIP])
    hosts.append(nodeIP)
    print("Added to Cassandra cluster")

def removeCassandraNode(nodeIP):
    print("----------------\nRemoving {} from Cassandra cluster\n-----------------".format(nodeIP))
    hostId = getHostId(nodeIP)
    subprocess.call(['cassandra/stopCluster.sh', nodeIP])
    subprocess.call('docker exec -ti cassandra-node nodetool removenode {}'.format(hostId).split(" "))
    print("Removed from Cassandra cluster")

def buildImages():
    subprocess.call(['./build.sh'])

# HOSTS
def createSwarm():
    global join_swarm_command
    try:
        byteOutput = subprocess.check_output('docker swarm init --advertise-addr {}'.format(manager).split(" "))
        join_swarm_command = byteOutput.decode('UTF-8').rstrip().splitlines()[4].strip()
        swarms.append(manager)
    except subprocess.CalledProcessError as e:
        print("Error in ls -a:\n", e.output)

def addToSwarm(nodeIP):
    print("----------------\nAdding {} to swarm\n-----------------".format(nodeIP))
    subprocess.call(['cassandra/create_volumes.sh', nodeIP])
    subprocess.call('sshpass -p hhhhiotwwg ssh student@{} {}'.format(nodeIP, join_swarm_command).split(" "))
    print("Added to swarm")
    swarms.append(nodeIP)

def removeFromSwarm(nodeIP):
    print("----------------\nRemoving {} from swarm\n-----------------".format(nodeIP))
    subprocess.call(['./removeFromSwarm.sh', nodeIP])
    print("Removed from swarm")

def removeSwarm():
    subprocess.call(['./removeFromSwarm.sh', *swarms])

def scaleService(service, replicas):
    subprocess.call('docker service scale {}={}'.format(service, replicas).split(" "))

# SYSTEM
def startSystem():
    retcode = subprocess.call('docker stack deploy -c docker-compose.yml url_stack --with-registry-auth'.split(" "))
    if retcode == 0:
        print("Started system")
        return
    print("Failed to start stack")

def stopSystem():
    retcode = subprocess.call('docker stack rm url_stack'.split(" "))
    if retcode == 0:
        print("Success")
        return
    print("Failed to stop system")

def setup():
    print("----------------\nStarting Cassandra\n-----------------")
    startCassandra()
    print("----------------\nFinished Cassandra\n-----------------\nBuilding and Publishing images")
    buildImages()
    print("----------------\nFinished publishing images\n-----------------\nStarting system")
    createSwarm()
    addToSwarm(hosts[1])
    addToSwarm(hosts[2])
    startSystem()

def stop():
    print("-----------------\nStopping system\n-----------------")
    stopSystem()
    removeCassandra()
    removeSwarm()

def showOptions():
    print("\nWelcome to the help option")
    print("help:")
    print("     - cassa <NODE_IP> : command to add a new node to the Cassandra cluster")
    print("     - cassr <NODE_IP> : command to remove a node from the Cassandra cluster")
    print("     - swarma <NODE_IP>: command to add a new machine to the existing swarm")
    print("     - swarmr <NODE_IP>: command to remove a machine from the existing swarm")
    print("     - scaleser <SERVICE_NAME> <NODE_IP>: command to scale service like url_stack_redis-secondary 4. Refer to localhost:8080")
    print("     - quit : to quit the system")

def processInput(inp: str):
    splitInp = inp.split(" ")
    if len(splitInp) < 2 or len(splitInp) > 3:
        print("consult the help guide on the right arguments")
        return
    cmd = splitInp[0]
    if cmd not in ACTION_TABLE.keys():
        print("not a supported option")
        return
    if len(splitInp) == 3:
        ACTION_TABLE[cmd](splitInp[1], splitInp[2])
        return
    nodeIP = splitInp[1]
    ACTION_TABLE[cmd](nodeIP)

ACTION_TABLE = {
    "cassa": addCassandraNode,
    "cassr": removeCassandraNode,
    "swarma": addToSwarm,
    "swarmr": removeFromSwarm,
    "scaleser": scaleService
}

def run_main():
    setup()
    print("Enter \"quit\" to stop the system. Enter \"help\" to know what options I offer.\nTo execute HTTP requests enter in another terminal")
    while True:
        inp = input("\nEnter command to execute: ")
        if inp == "quit":
            break
        elif inp == "help":
            showOptions()
        else:
            processInput(inp)
    stop()

if __name__ == "__main__":
    run_main()