#!/bin/bash
# Generate self-signed certificate for testing HTTPS locally

set -e

echo "🔐 Generating self-signed SSL certificate for local HTTPS testing"
echo "=================================================================="
echo ""

# Create certs directory
mkdir -p certs

# Generate certificate
openssl req -x509 -newkey rsa:4096 -nodes \
  -keyout certs/key.pem \
  -out certs/cert.pem \
  -days 365 \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=api.grandexchange.gg"

echo ""
echo "✅ Certificate generated!"
echo ""
echo "📁 Files created:"
echo "   - certs/cert.pem (certificate)"
echo "   - certs/key.pem (private key)"
echo ""
echo "🧪 To test locally with HTTPS:"
echo ""
echo "   export SSL_CERT_FILE=certs/cert.pem"
echo "   export SSL_KEY_FILE=certs/key.pem"
echo "   export PORT=8443"
echo "   go run main.go"
echo ""
echo "   curl -k https://localhost:8443/health"
echo ""
echo "⚠️  Note: This is for LOCAL TESTING ONLY"
echo "   For production, use Cloudflare Origin Certificates"
echo ""