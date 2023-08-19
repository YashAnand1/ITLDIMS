import csv
import json
import etcd3
import gspread
import pandas as pd
from oauth2client.service_account import ServiceAccountCredentials
from flask import Flask, request, jsonify
import argparse

app = Flask(__name__)

# Google Sheets credentials file path
credentials_file = '/home/additi/Downloads/additi-a2624591d964.json'

# Spreadsheet ID and sheet name
spreadsheet_id = '11uP_Ia1RRrtq3CMPZoFqTankBORsHpZ2DgMdTF4hoNU'
sheet_name = 'sheet'

# File paths
excel_file = '/tmp/etcd.xlsx'
csv_file = '/tmp/etcd.csv'

def convert_excel_to_csv(excel_file, csv_file):
    # Read the Excel file
    df = pd.read_excel(excel_file)

    # Write the data to a CSV file
    df.to_csv(csv_file, index=False, quoting=csv.QUOTE_NONNUMERIC)

# Connect to etcd
etcd_client = etcd3.client(host='172.18.0.2', port=2379)

def process_google_sheets():
    try:
        # Authenticate with Google Sheets API
        scope = ['https://www.googleapis.com/auth/spreadsheets']
        credentials = ServiceAccountCredentials.from_json_keyfile_name(credentials_file, scope)
        client = gspread.authorize(credentials)

        # Open the Google Spreadsheet
        spreadsheet = client.open_by_key(spreadsheet_id)
        worksheet = spreadsheet.worksheet(sheet_name)

        # Get all values from the sheet
        data = worksheet.get_all_values()

        # Write the data to a CSV file
        with open(csv_file, 'w', newline='') as file:
            writer = csv.writer(file)
            writer.writerows(data)

    except FileNotFoundError:
        print("Credentials file not found. Exiting.")
        exit(1)

def process_local_excel():
    try:
        # Convert XLSX to CSV
        convert_excel_to_csv(excel_file, csv_file)

    except FileNotFoundError:
        print("Excel file not found. Exiting.")
        exit(1)

def read_csv_data():
    # Read server details from CSV and store in etcd
    with open(csv_file, 'r') as file:
        csv_data = csv.reader(file)
        next(csv_data)  # Skip header row
        for row in csv_data:
            server_ip = row[0]
            server_type = row[1]
            hostname = row[2]
            cpu = row[3]
            ram = row[4]
            internal_disk = row[5]
            external_disk = row[6]
            internal_partition = row[7]
            external_partition = row[8]
            application = row[9]
            api = row[10]
            lvm = row[11]
            vg = row[12]
            pv = row[13]
            nfs = row[14]
            netmask = row[15]
            os = row[16]
            gateway = row[17]
            environment = row[18]

            # Create server data dictionary
            server_data = {
                'Type': server_type,
                'Hostname': hostname,
                'CPU': cpu,
                'RAM': ram,
                'Internal_Disk': internal_disk,
                'External_Disk': external_disk,
                'Internal_Partition': internal_partition,
                'External_Partition': external_partition,
                'Application': application,
                'API': api,
                'LVM': lvm,
                'VG': vg,
                'PV': pv,
                'NFS': nfs,
                'Netmask': netmask,
                'OS': os,
                'Gateway': gateway,
                'Environment': environment,
            }

            # Set key-value pairs in etcd
            etcd_key_type = f'/servers/{server_type}/{server_ip}'
            etcd_key_hostname = f'/servers/{server_type}/{server_ip}/hostname'
            etcd_key_cpu = f'/servers/{server_type}/{server_ip}/cpu'
            etcd_key_ram = f'/servers/{server_type}/{server_ip}/ram'
            etcd_key_internal_disk = f'/servers/{server_type}/{server_ip}/internal_disk'
            etcd_key_external_disk = f'/servers/{server_type}/{server_ip}/external_disk'
            etcd_key_internal_partition = f'/servers/{server_type}/{server_ip}/internal_partition'
            etcd_key_external_partition = f'/servers/{server_type}/{server_ip}/external_partition'
            etcd_key_application = f'/servers/{server_type}/{server_ip}/application'
            etcd_key_api = f'/servers/{server_type}/{server_ip}/api'
            etcd_key_lvm = f'/servers/{server_type}/{server_ip}/lvm'
            etcd_key_vg = f'/servers/{server_type}/{server_ip}/vg'
            etcd_key_pv = f'/servers/{server_type}/{server_ip}/pv'
            etcd_key_nfs = f'/servers/{server_type}/{server_ip}/nfs'
            etcd_key_os = f'/servers/{server_type}/{server_ip}/os'
            etcd_key_gateway = f'/servers/{server_type}/{server_ip}/gateway'
            etcd_key_netmask = f'/servers/{server_type}/{server_ip}/netmask'
            etcd_key_environment = f'/servers/{server_type}/{server_ip}/environment'
            etcd_value = json.dumps(server_data)

            etcd_client.put(etcd_key_type, etcd_value)
            etcd_client.put(etcd_key_hostname, hostname)
            etcd_client.put(etcd_key_cpu, cpu)
            etcd_client.put(etcd_key_ram, ram)
            etcd_client.put(etcd_key_internal_disk, internal_disk)
            etcd_client.put(etcd_key_external_disk, external_disk)
            etcd_client.put(etcd_key_internal_partition, internal_partition)
            etcd_client.put(etcd_key_external_partition, external_partition)
            etcd_client.put(etcd_key_application, application)
            etcd_client.put(etcd_key_api, api)
            etcd_client.put(etcd_key_lvm, lvm)
            etcd_client.put(etcd_key_vg, vg)
            etcd_client.put(etcd_key_pv, pv)
            etcd_client.put(etcd_key_nfs, nfs)
            etcd_client.put(etcd_key_netmask, netmask)
            etcd_client.put(etcd_key_os, os)
            etcd_client.put(etcd_key_gateway, gateway)
            etcd_client.put(etcd_key_environment, environment)

            # Read the key from etcd
            get_response = etcd_client.get(etcd_key_type)
            if get_response is not None:
                stored_data = json.loads(get_response[0])
                print("Read from etcd for server IP", server_ip, ":", stored_data)
            else:
                print("Key not found in etcd for server IP", server_ip)

        print("Server details added to etcd successfully.")

# GET API Route
@app.route('/servers/<server_type>/<server_ip>/<key>', methods=['GET'])
def get_server_key(server_type, server_ip, key):
    etcd_key = f'/servers/{server_type}/{server_ip}/{key}'
    get_response = etcd_client.get(etcd_key)
    if get_response is not None:
        value = get_response[0].decode('utf-8')
        return jsonify({key: value})
    else:
        return jsonify({'error': f'Key "{key}" not found.'}), 404


# POST API Route
@app.route('/servers/<server_type>/<server_ip>/<key>', methods=['POST'])
def post_server_key(server_type, server_ip, key):
    etcd_key = f'/servers/{server_type}/{server_ip}/{key}'
    value = request.get_json(force=True)
    if value is not None:
        etcd_client.put(etcd_key, value)
        return jsonify({etcd_key: value}), 200
    else:
        return jsonify({'error': 'Invalid JSON data'}), 400


# PUT API Route
@app.route('/servers/<server_type>/<server_ip>/<key>', methods=['PUT'])
def put_server_key(server_type, server_ip, key):
    etcd_key = f'/servers/{server_type}/{server_ip}/{key}'
    value = request.get_json(force=True)
    if value is not None:
        etcd_client.put(etcd_key, value)
        return jsonify({etcd_key: value}), 200
    else:
        return jsonify({'error': 'Invalid JSON data'}), 400

# DELETE API Route
@app.route('/servers/<server_type>/<server_ip>/<key>', methods=['DELETE'])
def delete_server_key(server_type, server_ip, key):
    etcd_key = f'/servers/{server_type}/{server_ip}/{key}'
    deleted = etcd_client.delete(etcd_key)
    if deleted:
        return jsonify({'message': 'Value deleted successfully'}), 200
    else:
        return jsonify({'error': 'Server not found'}), 404



if __name__ == '__main__':
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Process server data and store in etcd')
    parser.add_argument('--google-sheets', action='store_true', help='Process data from Google Sheets')
    parser.add_argument('--local-excel', action='store_true', help='Process data from local Excel file')
    args = parser.parse_args()

    if args.google_sheets:
        process_google_sheets()
    elif args.local_excel:
        process_local_excel()
    else:
        print("No input option specified. Please use either --google-sheets or --local-excel.")
        exit(1)

    read_csv_data()

    # Run the Flask app
    app.run(host='0.0.0.0', port=5000)

