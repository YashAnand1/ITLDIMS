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
excel_file = '/home/additi/etcd.xlsx'
csv_file = '/home/additi/etcd.csv'

def convert_excel_to_csv(excel_file, csv_file):
    # Read the Excel file
    df = pd.read_excel(excel_file)

    # Write the data to a CSV file
    df.to_csv(csv_file, index=False, quoting=csv.QUOTE_NONNUMERIC)

# Connect to etcd
etcd_client = etcd3.client(host='172.17.0.2', port=2379)

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
        headers = next(csv_data)  # Get header row
        for row in csv_data:
            server_ip = row[0]
            server_type = row[1]
            server_data = {}

            # Create server data dictionary
            for i in range(2, len(headers)):
                header = headers[i]
                value = row[i]
                server_data[header] = value

            # Set key-value pairs in etcd for each data field
            for header, value in server_data.items():
                etcd_key = f'/servers/{server_type}/{server_ip}/{header}'
                etcd_value = json.dumps(value)
                etcd_client.put(etcd_key, etcd_value)

            # Set key-value pair for server data
            etcd_key_data = f'/servers/{server_type}/{server_ip}/data'
            etcd_value_data = json.dumps(server_data)
            etcd_client.put(etcd_key_data, etcd_value_data)

        print("Server details added to etcd successfully.")


# GET API Route
@app.route('/servers/<server_type>/<server_ip>/<key>', methods=['GET'])
def get_server_key(server_type, server_ip, key):
    etcd_key = f'/servers/{server_type}/{server_ip}/{key}'

    # Retrieve the current value and its version
    current_value, current_metadata = etcd_client.get(etcd_key)

    if current_value is not None:
        # Extract the version from the metadata
        current_version = current_metadata.version

        # Retrieve the previous value and its version
        previous_value = None
        previous_version = current_version - 1
        if previous_version >= 0:
            revisions = etcd_client.get_prefix(etcd_key, sort_order='descend', sort_target='create', limit=previous_version + 1)
            revisions_list = list(revisions)
            if len(revisions_list) > 1:
                previous_value = revisions_list[1][0].decode('utf-8')  # Retrieve the previous value correctly


        # Convert the current value to a string
        current_value = current_value.decode('utf-8')

        # Create the response JSON
        response = {
            'key': key,
            'current_value': current_value,
            'current_version': current_version,
            'previous_value': previous_value,
            'previous_version': previous_version
        }
        return jsonify(response)
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

