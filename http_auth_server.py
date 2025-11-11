#!/usr/bin/env python3
"""Simple HTTP Basic Auth server for testing bolt-proxy authentication"""

from http.server import BaseHTTPRequestHandler, HTTPServer
import base64

class AuthHandler(BaseHTTPRequestHandler):
    """HTTP handler that validates Basic Auth credentials"""

    def do_GET(self):
        """Handle GET requests with Basic Auth"""
        auth_header = self.headers.get('Authorization')

        if auth_header and auth_header.startswith('Basic '):
            try:
                # Decode the base64 credentials
                credentials = base64.b64decode(auth_header[6:]).decode('utf-8')
                username, password = credentials.split(':', 1)

                # Validate credentials (neo4j/password)
                if username == 'neo4j' and password == 'password':
                    self.send_response(200)
                    self.send_header('Content-Type', 'text/plain')
                    self.end_headers()
                    self.wfile.write(b'Authenticated')
                    print(f'AUTH SUCCESS: {username}')
                    return
                else:
                    print(f'AUTH FAILED: Invalid credentials for {username}')
            except Exception as e:
                print(f'AUTH ERROR: {e}')

        # Unauthorized
        self.send_response(401)
        self.send_header('WWW-Authenticate', 'Basic realm="Bolt Proxy Auth"')
        self.send_header('Content-Type', 'text/plain')
        self.end_headers()
        self.wfile.write(b'Unauthorized')
        print('AUTH FAILED: No valid authorization header')

    def log_message(self, format, *args):
        """Custom log message format"""
        print(f'[AUTH] {args[0]} - {args[1]}')

if __name__ == '__main__':
    server_address = ('0.0.0.0', 8081)
    httpd = HTTPServer(server_address, AuthHandler)
    print('Starting HTTP Basic Auth server on port 8081...')
    print('Valid credentials: neo4j/password')
    httpd.serve_forever()
