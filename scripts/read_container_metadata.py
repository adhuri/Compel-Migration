import simplejson as json
import os
import optparse


parser = optparse.OptionParser()
parser.add_option('-n', '--name', dest='name', help='Container Name/Id')
(options, args) = parser.parse_args()
containerName = options.name


#docker_json = os.system('docker inspect ' + containerName)

import subprocess

proc = subprocess.Popen(["docker inspect "+containerName], stdout=subprocess.PIPE, shell=True)
(json_string, err) = proc.communicate()
print json_string[1:-2]


parsed_json = json.loads(json_string[4:-3])

#print(parsed_json['Config']['Env'])
if parsed_json['Config']['Env'] != None:
    for k in parsed_json['Config']['Env']:
        k = k.replace("=", " ")
        print "ENV " + k
#print(parsed_json['Config']['Cmd'])
if parsed_json['Config']['Cmd'] != None:
    for k in parsed_json['Config']['Cmd']:
        #k = k.replace("=", " ")
        print "CMD " + k
print(parsed_json['Config']['ExposedPorts'])
print(parsed_json['Config']['Entrypoint'])
print(parsed_json['NetworkSettings']['Ports'])
