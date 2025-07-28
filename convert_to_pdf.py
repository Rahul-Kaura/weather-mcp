#!/usr/bin/env python3
"""
Convert HTML documentation to PDF
Requires: pip install weasyprint
"""

import os
import sys
from pathlib import Path

def convert_html_to_pdf():
    """Convert the HTML documentation to PDF"""
    try:
        import weasyprint
        
        # Check if HTML file exists
        html_file = "generate_pdf.html"
        if not os.path.exists(html_file):
            print(f"❌ HTML file '{html_file}' not found!")
            return False
        
        # Convert to PDF
        print("🔄 Converting HTML to PDF...")
        weasyprint.HTML(filename=html_file).write_pdf("Weather_MCP_Implementation_Guide.pdf")
        
        # Check if PDF was created
        if os.path.exists("Weather_MCP_Implementation_Guide.pdf"):
            file_size = os.path.getsize("Weather_MCP_Implementation_Guide.pdf")
            print(f"✅ PDF created successfully!")
            print(f"📄 File: Weather_MCP_Implementation_Guide.pdf")
            print(f"📏 Size: {file_size:,} bytes ({file_size/1024:.1f} KB)")
            return True
        else:
            print("❌ PDF creation failed!")
            return False
            
    except ImportError:
        print("❌ weasyprint not installed!")
        print("💡 Install with: pip install weasyprint")
        return False
    except Exception as e:
        print(f"❌ Error converting to PDF: {e}")
        return False

def main():
    """Main function"""
    print("🌤️ Weather MCP Implementation Guide - PDF Generator")
    print("=" * 60)
    
    # Check if weasyprint is available
    try:
        import weasyprint
        print("✅ weasyprint is available")
    except ImportError:
        print("❌ weasyprint not found!")
        print("💡 Please install it with: pip install weasyprint")
        print("💡 Or use the Markdown file directly: WEATHER_MCP_IMPLEMENTATION_GUIDE.md")
        return
    
    # Convert HTML to PDF
    success = convert_html_to_pdf()
    
    if success:
        print("\n🎉 Success! You now have:")
        print("📄 Weather_MCP_Implementation_Guide.pdf - PDF version")
        print("📝 WEATHER_MCP_IMPLEMENTATION_GUIDE.md - Markdown version")
        print("\n💡 Both files are ready for your meeting!")
    else:
        print("\n💡 You can still use the Markdown file:")
        print("📝 WEATHER_MCP_IMPLEMENTATION_GUIDE.md")

if __name__ == "__main__":
    main() 