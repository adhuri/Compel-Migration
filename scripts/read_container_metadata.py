import os
import optparse
import importlib
import subprocess


## Parsing Command-Line arguments
parser = optparse.OptionParser()
parser.add_option('-n', '--name', dest='name', help='Container Name/Id')
parser.add_option('-u', '--userame', dest='username', help='User Name on Host')
(options, args) = parser.parse_args()
containerName = options.name
userName = options.username


## Importing or Installing simplejson module to make sure python code doesn't fail
try:
    importlib.import_module("simplejson")
except ImportError:
    import pip
    pip.main(['install', "simplejson"])
finally:
    globals()["simplejson"] = importlib.import_module("simplejson")


## Execute docker inspect command for a given container
proc = subprocess.Popen(["docker inspect "+containerName], stdout=subprocess.PIPE, shell=True)
(json_string, err) = proc.communicate()


## Parse the returned JSON into a dict
parsed_json = simplejson.loads(json_string[4:-3])


## Open a file for output
f = open('/home/'+userName+'/'+containerName+'_metadata.conf', 'w')


## Parse ENV variables
#print(parsed_json['Config']['Env'])
if parsed_json['Config']['Env'] != None:
    for k in parsed_json['Config']['Env']:
        k = k.replace("=", " ")
        f.write(("ENV " + k+"\n").encode("utf-8"))


## Parse CMD command
#print(parsed_json['Config']['Cmd'])
if parsed_json['Config']['Cmd'] != None:
    for k in parsed_json['Config']['Cmd']:
        f.write("CMD " + k+"\n")


## Parse EXPOSED Ports
#print(parsed_json['Config']['ExposedPorts'])
if parsed_json['Config']['ExposedPorts'] != None :
    for k in parsed_json['Config']['ExposedPorts']:
        f.write("EXPOSE " + k+"\n")


## Parse ENTRYPOINT for container
#print(parsed_json['Config']['Entrypoint'])
if parsed_json['Config']['Entrypoint'] != None :
    for k in parsed_json['Config']['Entrypoint']:
        f.write("ENTRYPOINT " + k+"\n")


## Parse Port Mapping for container
#print(parsed_json['NetworkSettings']['Ports'])
if parsed_json['NetworkSettings']['Ports'] != None:
    for k in parsed_json['NetworkSettings']['Ports']:
        if parsed_json['NetworkSettings']['Ports'][k] != None:
            f.write("-p "+parsed_json['NetworkSettings']['Ports'][k][0]['HostPort']+":"+k[0:k.index("/")]+"\n")



## Close output file
f.close()  # you can omit in most cases as the destructor will call it
