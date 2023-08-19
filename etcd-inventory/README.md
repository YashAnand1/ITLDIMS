# Etcd Server Configuration

This script connects to an Etcd server and performs server configuration tasks using data from a Google Sheet or an Excel file. The server details are stored in Etcd, and an API is provided to retrieve, update, and delete server configuration values.

## Prerequisites

Before running the script, ensure that you have the following installed:

- Python 3
- `etcd3` Python library
- `gspread` Python library
- `pandas` Python library
- `oauth2client` Python library
- `openpyxl` Python library
- `Flask` Python library

## Setup

1. Clone the repository and navigate to the project directory:

```shell
git clone https://github.com/Keen-And-Able/etcd-inventory.git
cd etcd-inventory
```

2. Install the required Python libraries:

```shell
pip install etcd3 gspread pandas oauth2client openpyxl Flask
```

## Usage

The script can process server data from either Google Sheets or a local Excel file. You can choose the data source using command-line arguments.

### Google Sheets

To process data from Google Sheets, follow these steps:

1. Create a service account and download the JSON credentials file. See the [Google Sheets API documentation](https://developers.google.com/sheets/api/guides/authorizing#creating_a_service_account) for detailed instructions.

2. Move the downloaded JSON credentials file to the project directory and set the `credentials_file` variable in the script to the file path.

3. Set the `spreadsheet_id` variable in the script to the ID of the Google Spreadsheet containing the server data.

4. Set the `sheet_name` variable in the script to the name of the sheet containing the server data.

5. Run the script with the `--google-sheets` option:

   ```shell
   python3 etcd.py --google-sheets
   ```

### Local Excel File

To process data from a local Excel file, follow these steps:

1. Place your Excel file in the project directory.

2. Set the `excel_file` variable in the script to the file path of your Excel file.

3. Run the script with the `--local-excel` option:

   ```shell
   python3 etcd.py --local-excel
   ```


## API Routes

The script provides the following API routes:

- `GET /servers/<server_type>/<server_ip>/<key>`: Retrieves the value of the specified key for a given server.
- `POST /servers/<server_type>/<server_ip>/<key>`: Sets the value of the specified key for a given server.
- `PUT /servers/<server_type>/<server_ip>/<key>`: Updates the value of the specified key for a given server.
- `DELETE /servers/<server_type>/<server_ip>/<key>`: Deletes the specified key for a given server.

Note:Replace `<server_type>`, `<server_ip>`, and `<key>` with the appropriate values.


---

