/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Don't advertise the framework/version.
  poweredByHeader: false,

  async headers() {
    return [
      {
        source: "/:path*",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "X-Frame-Options", value: "SAMEORIGIN" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
          // NOTE: no restrictive Permissions-Policy here — the proctored exam
          // needs camera + getDisplayMedia, which a strict policy would block.
        ],
      },
    ];
  },
};

module.exports = nextConfig;
