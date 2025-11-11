#!/usr/bin/env python3
"""Test script for bolt-proxy authentication and CRUD operations"""

import time
import sys
import os

try:
    from neo4j import GraphDatabase
except ImportError:
    print("Error: neo4j driver not installed. Install with: pip install neo4j")
    sys.exit(1)


class BoltProxyTester:
    """Test client for bolt-proxy authentication and CRUD operations"""

    def __init__(self, uri="bolt://bolt_proxy:7474", correct_user="neo4j", correct_password="password"):
        self.uri = uri
        self.correct_user = correct_user
        self.correct_password = correct_password
        self.driver = None

    def close(self):
        """Close the driver connection"""
        if self.driver:
            self.driver.close()

    def test_auth_success(self):
        """Test successful authentication"""
        print("\n=== Testing Successful Authentication ===")
        try:
            self.driver = GraphDatabase.driver(
                self.uri,
                auth=(self.correct_user, self.correct_password)
            )
            # Verify connection works
            with self.driver.session() as session:
                result = session.run("RETURN 1 as num")
                record = result.single()
                print(f"✓ Auth SUCCESS: Connected as {self.correct_user}")
                return True
        except Exception as e:
            print(f"✗ Auth FAILED: {e}")
            return False

    def test_auth_failure(self):
        """Test failed authentication"""
        print("\n=== Testing Failed Authentication ===")
        try:
            bad_driver = GraphDatabase.driver(
                self.uri,
                auth=("wrong_user", "wrong_password")
            )
            with bad_driver.session() as session:
                session.run("RETURN 1")
            bad_driver.close()
            print("✗ Auth should have FAILED but didn't")
            return False
        except Exception as e:
            print(f"✓ Auth FAILED as expected: {type(e).__name__}")
            return True

    def test_create_operation(self):
        """Test CREATE operation"""
        print("\n=== Testing CREATE Operation ===")
        try:
            with self.driver.session() as session:
                result = session.run(
                    "CREATE (p:Person {name: $name, age: $age}) RETURN p",
                    name="Alice", age=30
                )
                record = result.single()
                print(f"✓ CREATE: Created person node")
                return True
        except Exception as e:
            print(f"✗ CREATE failed: {e}")
            return False

    def test_read_operation(self):
        """Test READ operation"""
        print("\n=== Testing READ Operation ===")
        try:
            with self.driver.session() as session:
                result = session.run("MATCH (p:Person) RETURN p.name, p.age LIMIT 5")
                count = 0
                for record in result:
                    count += 1
                    print(f"  - Found: {record['p.name']}, age {record['p.age']}")
                print(f"✓ READ: Found {count} person(s)")
                return True
        except Exception as e:
            print(f"✗ READ failed: {e}")
            return False

    def test_update_operation(self):
        """Test UPDATE operation"""
        print("\n=== Testing UPDATE Operation ===")
        try:
            with self.driver.session() as session:
                result = session.run(
                    "MATCH (p:Person {name: 'Alice'}) SET p.age = $new_age RETURN p",
                    new_age=31
                )
                if result.single():
                    print(f"✓ UPDATE: Updated Alice's age")
                    return True
                else:
                    print(f"✗ UPDATE: No nodes updated")
                    return False
        except Exception as e:
            print(f"✗ UPDATE failed: {e}")
            return False

    def test_delete_operation(self):
        """Test DELETE operation"""
        print("\n=== Testing DELETE Operation ===")
        try:
            with self.driver.session() as session:
                result = session.run(
                    "MATCH (p:Person {name: 'Alice'}) DELETE p"
                )
                result.consume()
                print(f"✓ DELETE: Deleted person node")
                return True
        except Exception as e:
            print(f"✗ DELETE failed: {e}")
            return False

    def cleanup(self):
        """Clean up test data"""
        print("\n=== Cleanup ===")
        try:
            with self.driver.session() as session:
                session.run("MATCH (p:Person) DELETE p")
                print("✓ Cleaned up all Person nodes")
        except Exception as e:
            print(f"✗ Cleanup failed: {e}")

    def run_all_tests(self, iterations=1):
        """Run all tests for specified iterations"""
        print(f"\n{'='*60}")
        print(f"Bolt Proxy Authentication & CRUD Operation Tester")
        print(f"Target: {self.uri}")
        print(f"Iterations: {iterations}")
        print(f"{'='*60}")

        for i in range(iterations):
            if iterations > 1:
                print(f"\n\n{'#'*60}")
                print(f"# Iteration {i+1}/{iterations}")
                print(f"{'#'*60}")

            # Test authentication
            auth_success = self.test_auth_success()
            if not auth_success:
                print("\n✗ Cannot proceed without successful authentication")
                return

            self.test_auth_failure()

            # Test CRUD operations
            self.test_create_operation()
            time.sleep(0.1)

            self.test_read_operation()
            time.sleep(0.1)

            self.test_update_operation()
            time.sleep(0.1)

            self.test_read_operation()  # Read again to verify update
            time.sleep(0.1)

            self.test_delete_operation()
            time.sleep(0.1)

            if i < iterations - 1:
                time.sleep(1)  # Pause between iterations

        # Final cleanup
        self.cleanup()
        self.close()

        print(f"\n{'='*60}")
        print("Testing Complete!")
        print(f"{'='*60}\n")


if __name__ == '__main__':
    # Read configuration from environment variables
    bolt_host = os.getenv("BOLT_HOST", "localhost")
    bolt_port = os.getenv("BOLT_PORT", "8888")
    bolt_uri = os.getenv("BOLT_URI", f"bolt://{bolt_host}:{bolt_port}")
    bolt_user = os.getenv("BOLT_USER", "neo4j")
    bolt_password = os.getenv("BOLT_PASSWORD", "password")

    # Parse command line arguments
    iterations = 1

    if len(sys.argv) > 1:
        try:
            iterations = int(sys.argv[1])
        except ValueError:
            print(f"Usage: {sys.argv[0]} [iterations]")
            print(f"")
            print(f"Environment variables:")
            print(f"  BOLT_URI      - Full bolt URI (default: bolt://localhost:8888)")
            print(f"  BOLT_HOST     - Bolt host (default: localhost)")
            print(f"  BOLT_PORT     - Bolt port (default: 8080)")
            print(f"  BOLT_USER     - Username (default: neo4j)")
            print(f"  BOLT_PASSWORD - Password (default: password)")
            sys.exit(1)

    # Run tests
    tester = BoltProxyTester(uri=bolt_uri, correct_user=bolt_user, correct_password=bolt_password)
    try:
        tester.run_all_tests(iterations=iterations)
    except KeyboardInterrupt:
        print("\n\nInterrupted by user")
        tester.cleanup()
        tester.close()
    except Exception as e:
        print(f"\n\nUnexpected error: {e}")
        tester.close()
        sys.exit(1)
