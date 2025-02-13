import ssl
from _hashlib import get_fips_mode
import hashlib
import unittest


class TestCompliance(unittest.TestCase):
    def test_openssl_version(self):
        self.assertEqual(ssl.OPENSSL_VERSION, "OpenSSL 3.0.9 30 May 2023")

    def test_fips_mode(self):
        self.assertEqual(get_fips_mode(), 1)

    def test_hashes(self):
        data = b"foo"

        self.assertIsNotNone(hashlib.sha1(data))
        self.assertIsNotNone(hashlib.new("sha1", data))
        self.assertIsNotNone(hashlib.sha224(data))
        self.assertIsNotNone(hashlib.new("sha224", data))
        self.assertIsNotNone(hashlib.sha256(data))
        self.assertIsNotNone(hashlib.new("sha256", data))
        self.assertIsNotNone(hashlib.sha384(data))
        self.assertIsNotNone(hashlib.new("sha384", data))
        self.assertIsNotNone(hashlib.sha3_224(data))
        self.assertIsNotNone(hashlib.new("sha3_224", data))
        self.assertIsNotNone(hashlib.sha3_256(data))
        self.assertIsNotNone(hashlib.new("sha3_256", data))
        self.assertIsNotNone(hashlib.sha3_384(data))
        self.assertIsNotNone(hashlib.new("sha3_384", data))
        self.assertIsNotNone(hashlib.sha3_512(data))
        self.assertIsNotNone(hashlib.new("sha3_512", data))
        self.assertIsNotNone(hashlib.sha512(data))
        self.assertIsNotNone(hashlib.new("sha512", data))

        with self.assertRaises(ValueError):
            hashlib.md5(data)
        with self.assertRaises(ValueError):
            hashlib.new("md5", data)
        with self.assertRaises(ValueError):
            hashlib.blake2b(data)
        with self.assertRaises(ValueError):
            hashlib.new("blake2b", data)
        with self.assertRaises(ValueError):
            hashlib.blake2s(data)
        with self.assertRaises(ValueError):
            hashlib.new("blake2s", data)


if __name__ == "__main__":
    unittest.main()
