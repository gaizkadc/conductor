#!/usr/bin/python

'''
This script is a temporal workaround to test the integration between
- system model
- conductor
- asm client

The script loads a bunch of data entries into the system model using budo, upload an influxdb image to a colonize
deployed instance and then uses conductor CLI to start the deployment of an application.

Launch the script with something like:

./integration.py --basedir=/Users/juan/daisho/daisho_workspace/src/github.com/daishogroup --appmgr=/Users/juan/daisho/appmgr/bazel-bin/asmcli/asmcli --desc=/Users/juan/daisho/appdevkit/bazel-out/local-fastbuild/bin/packages/influxdb/influxdb-asm-package.tar.gz


The script assumes you have a working appmgr (probably deployed with colony) and it points to 172.28.128.4

'''



import subprocess
import json
import argparse
import sys



DESCRIPTION=""" 
This script is a temporal workaround to test the integration between
- system model
- conductor
- asm client

The script loads a bunch of data entries into the system model using budo, upload an influxdb image to a colonize
deployed instance and then uses conductor CLI to start the deployment of an application.

Launch the script with something like:

./integration.py --basedir=/Users/juan/daisho/daisho_workspace/src/github.com/daishogroup --appmgr=/Users/juan/daisho/appmgr/bazel-bin/asmcli/asmcli --desc=/Users/juan/daisho/appdevkit/bazel-out/local-fastbuild/bin/packages/influxdb/influxdb-asm-package.tar.gz


The script assumes you have a working appmgr (probably deployed with colony) and it points to 172.28.128.4
"""

BASE_DIR=None

# Where is the system model client
SM_CLI=None
SM_IP='localhost'
SM_PORT='88000'

# Where is the conductor client
CONDUCTOR_CLI=None
CONDUCTOR_IP='localhost'
CONDUCTOR_PORT='90000'


# Where is the asmcli to pre-upload an application package
APPMGR_CLI=''
DESCRIPTOR_PATH=''

# The ip of the node to deploy
NODE1_IP = '172.28.128.4'
NODE2_IP = '172.28.128.5'
NODE3_IP = '172.28.128.6'

USERNAME='user'


def run_command(args,expected_json=True):
    '''
    Generic function to launch commands.
    :param args:
    :return:
    '''
    process = subprocess.Popen(args, stdout=subprocess.PIPE)
    print("--> {}".format(process.args))
    result, err = process.communicate()
    if err != None:
        print(err)
        return

    if expected_json:
        json_object = json.loads(result)
        print(json.dumps(json_object, indent=4, sort_keys=True))
        return json_object
    else:
        print(result.decode('utf-8'))

def run_sm(arguments):
    '''
    Run a system model command
    :param args:
    :return:
    '''
    elements = [SM_CLI,"--ip=%s"%SM_IP, "--port=%d"%SM_PORT]
    elements.extend(arguments)
    return run_command(elements)

def run_conductor(arguments):
    '''
    Run a conductor command
    :param args:
    :return:
    '''
    elements = [CONDUCTOR_CLI,"--ip=%s"%CONDUCTOR_IP, "--port=%d"%CONDUCTOR_PORT]
    elements.extend(arguments)
    return run_command(elements)

def start_clusters():
    '''
    Launch clusters for the testing
    :return:
    '''
    # TODO

def start_services():
    '''
    Start the services involved into the integration
    :return:
    '''
    # TODO


def package_upload():
    '''
    Upload the package to run
    :return:
    '''
    print('Upload package')
    run_command([APPMGR_CLI, '--ip=%s'%NODE1_IP, '--port=30088', 'package', 'upload', DESCRIPTOR_PATH],False)
    print('Check available manifests')
    run_command([APPMGR_CLI, '--ip=%s'%NODE1_IP, '--port=30088', 'manifest', 'list'], False)


def initialize_system_model():
    print('Initialize system model')
    print('Create initial network')
    network = run_sm(['network', 'add', 'network1'])

    print('Create clusters')
    cluster1 = run_sm(['cluster', 'add', network['id'], 'cluster1', 'cloud'])
    cluster1 = run_sm(['cluster', 'update', '--status=DEPLOYED', '--type=cloud', network['id'], cluster1['id']])
    cluster2 = run_sm(['cluster', 'add', network['id'], 'cluster2', 'gateway'])
    cluster2 = run_sm(['cluster', 'update', '--status=DEPLOYED', '--type=gateway', network['id'], cluster2['id']])
    cluster3 = run_sm(['cluster', 'add', network['id'], 'cluster3', 'edge'])
    cluster3 = run_sm(['cluster', 'update', '--status=DEPLOYED', '--type=edge', network['id'], cluster3['id']])
    cluster4 = run_sm(['cluster', 'add', network['id'], 'cluster4', 'edge'])
    cluster4 = run_sm(['cluster', 'update', '--status=DEPLOYED', '--type=edge', network['id'], cluster4['id']])


    print('Add nodes')
    node1 = run_sm(['node', 'add', network['id'], cluster1['id'], 'node1',
                         NODE1_IP, '0.0.0.0', 'user','--installed'])
    node2 = run_sm(['node', 'add', network['id'], cluster2['id'], 'node2',
                         NODE2_IP, '0.0.0.0', 'user', '--installed'])
    node3 = run_sm(['node', 'add', network['id'], cluster3['id'], 'node3',
                         NODE3_IP, '0.0.0.0', 'user', '--installed'])
    node4 = run_sm(['node', 'add', network['id'], cluster4['id'], 'node4',
                    NODE3_IP, '0.0.0.0', 'user', '--installed'])


    print('Add descriptor')
    descriptor = run_sm(['application', 'descriptor', 'add', network['id'], 'influxdb', 'influxdb', '0.2.1',
                              'gateway', '8888'])

    print('Done')
    return {'network': network,
            'cluster1': cluster1,
            'cluster2': cluster2,
            'cluster3': cluster3,
            'cluster4': cluster4,
            'node1': node1,
            'node2': node2,
            'node3': node3,
            'descriptor': descriptor}

def test_deploy(entities):

    print('Test deployment')
    result = run_command([CONDUCTOR_CLI, 'orchestrator', 'deploy', entities['network']['id'],
                          entities['descriptor']['id'], entities['descriptor']['label'], entities['descriptor']['name']])
    entities['instance']=result
    print('Tested')


def test_undeploy(entities):
    print('Test undeploy')
    result = run_command([CONDUCTOR_CLI, 'orchestrator', 'undeploy', entities['network']['id'],
                          entities['instance']['deployedId']],False)
    print('Tested')


def check_deployed(entities):
    print('Check deployed instance')
    result = run_command([APPMGR_CLI, '--ip=%s'%NODE1_IP, '--port=30088', 'app', 'list'], False)
    print('Done')

def stop_instance():
    print('Stop running instance')
    run_command([APPMGR_CLI, '--ip=%s'%NODE1_IP, '--port=30088', 'app', 'stop', 'influxdb'], False)

if __name__=='__main__':

    parser = argparse.ArgumentParser(description=DESCRIPTION,formatter_class=argparse.RawTextHelpFormatter)
    parser.add_argument('--basedir', dest='basedir', type=str, action='store', required=True,
                        help='Path where the daisho code is available')
    parser.add_argument('--appmgr', dest='appmgr', type=str, action='store', required=True,
                        help='Path where the compiled appmgr is available')
    parser.add_argument('--desc', dest='desc', action='store', type=str, required=True,
                        help='Path where the influxdb descriptor can be found')
    # options regarding the configuration of IPs and ports
    parser.add_argument('--sm_ip', dest='sm_ip', type=str, action='store', required=False, default='127.0.0.1',
                        help='System model IP')
    parser.add_argument('--sm_port', dest='sm_port', type=int, action='store', required=False, default=8800,
                        help='System model IP')
    parser.add_argument('--conductor_ip', dest='conductor_ip', type=str, action='store', required=False, default='127.0.0.1',
                        help='System model IP')
    parser.add_argument('--conductor_port', dest='conductor_port', type=int, action='store', required=False, default=9000,
                    help='Conductor port')


    args = vars(parser.parse_args(sys.argv[1:]))
    print(args)

    BASE_DIR=args['basedir']
    APPMGR_CLI=args['appmgr']
    DESCRIPTOR_PATH=args['desc']
    SM_IP=args['sm_ip']
    SM_PORT=args['sm_port']
    CONDUCTOR_IP=args['conductor_ip']
    CONDUCTOR_PORT=args['conductor_port']

    SM_CLI=BASE_DIR+'/system-model/bazel-bin/system-model-cli'
    CONDUCTOR_CLI=BASE_DIR+'/conductor/bazel-bin/conductor-cli'

    print('System model client located at: %s'%SM_CLI)
    print('Conductor client located at %s'%CONDUCTOR_CLI)

    package_upload()
    # Stop before running just for safety
    stop_instance()
    entities = initialize_system_model()
    # Deploy
    test_deploy(entities)
    # Check if it was already deployed
    check_deployed(entities)
    # Undeploy the previously deployed
    test_undeploy(entities)
    # Check there is nothing there
    check_deployed(entities)

