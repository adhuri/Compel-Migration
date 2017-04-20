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
#print json_string[1:-2]


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
#print(parsed_json['Config']['ExposedPorts'])
if parsed_json['Config']['ExposedPorts'] != None :
    for k in parsed_json['Config']['ExposedPorts']:
        #k = k.replace("=", " ")
        print "EXPOSE " + k
#print(parsed_json['Config']['Entrypoint'])
if parsed_json['Config']['Entrypoint'] != None :
    for k in parsed_json['Config']['Entrypoint']:
        #k = k.replace("=", " ")
        print "ENTRYPOINT " + k
#print(parsed_json['NetworkSettings']['Ports'])
if parsed_json['NetworkSettings']['Ports'] != None:
    for k in parsed_json['NetworkSettings']['Ports']:
        if parsed_json['NetworkSettings']['Ports'][k] != None:
            print "-p "+parsed_json['NetworkSettings']['Ports'][k][0]['HostPort']+":"+k[0:k.index("/")]
