#!/usr/bin/env python3
"""
Test script for bolt-proxy.

This script connects to the bolt-proxy and runs basic queries to verify
it's working correctly.

Requirements:
    pip install neo4j

Usage:
    python test_bolt_proxy.py [--host HOST] [--port PORT] [--user USER] [--password PASSWORD]
"""

import argparse
import sys
from neo4j import GraphDatabase


class BoltProxyTester:
    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))

    def close(self):
        self.driver.close()

    def test_connection(self):
        """Test basic connection to bolt-proxy."""
        print("Testing connection...")
        try:
            with self.driver.session() as session:
                result = session.run("RETURN 1 AS num")
                record = result.single()
                assert record["num"] == 1
                print("✓ Connection successful")
                return True
        except Exception as e:
            print(f"✗ Connection failed: {e}")
            return False

    def test_create_node(self):
        """Test creating a node."""
        print("\nTesting node creation...")
        try:
            with self.driver.session() as session:
                # Create a test node
                result = session.run(
                    "CREATE (n:TestNode {name: $name, timestamp: timestamp()}) RETURN n",
                    name="BoltProxyTest"
                )
                node = result.single()["n"]
                print(f"✓ Created node: {dict(node)}")
                return True
        except Exception as e:
            print(f"✗ Node creation failed: {e}")
            return False

    def test_query_node(self):
        """Test querying nodes."""
        print("\nTesting node query...")
        try:
            with self.driver.session() as session:
                result = session.run(
                    "MATCH (n:TestNode {name: $name}) RETURN n, id(n) AS id",
                    name="BoltProxyTest"
                )
                records = list(result)
                if records:
                    print(f"✓ Found {len(records)} node(s)")
                    for record in records:
                        print(f"  Node ID: {record['id']}, Properties: {dict(record['n'])}")
                    return True
                else:
                    print("✗ No nodes found")
                    return False
        except Exception as e:
            print(f"✗ Query failed: {e}")
            return False

    def test_cleanup(self):
        """Clean up test data."""
        print("\nCleaning up test data...")
        try:
            with self.driver.session() as session:
                result = session.run("MATCH (n:TestNode {name: $name}) DELETE n", name="BoltProxyTest")
                print("✓ Cleanup successful")
                return True
        except Exception as e:
            print(f"✗ Cleanup failed: {e}")
            return False

    def run_all_tests(self):
        """Run all tests."""
        print("=" * 60)
        print("Bolt Proxy Test Suite")
        print("=" * 60)

        tests = [
            self.test_connection,
            self.test_create_node,
            self.test_query_node,
            self.test_cleanup,
        ]

        passed = 0
        failed = 0

        for test in tests:
            if test():
                passed += 1
            else:
                failed += 1

        print("\n" + "=" * 60)
        print(f"Results: {passed} passed, {failed} failed")
        print("=" * 60)

        return failed == 0


def main():
    parser = argparse.ArgumentParser(description="Test bolt-proxy connection and functionality")
    parser.add_argument("--host", default="localhost", help="Bolt-proxy host (default: localhost)")
    parser.add_argument("--port", default="8080", help="Bolt-proxy port (default: 8080)")
    parser.add_argument("--user", default="", help="Username for authentication")
    parser.add_argument("--password", default="", help="Password for authentication")

    args = parser.parse_args()

    uri = f"bolt://{args.host}:{args.port}"
    print(f"Connecting to: {uri}")
    print(f"User: {args.user if args.user else '(none)'}\n")

    tester = BoltProxyTester(uri, args.user, args.password)

    try:
        success = tester.run_all_tests()
        sys.exit(0 if success else 1)
    finally:
        tester.close()


if __name__ == "__main__":
    main()
