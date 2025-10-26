"""
VideoDisco Phase 0 Setup Verification Test

This script verifies that all Phase 0 requirements are properly installed.
Run this to confirm your environment is ready for Phase 1.
"""

import sys
import importlib

def test_import(module_name, display_name=None):
    """Test if a module can be imported and print version info"""
    display_name = display_name or module_name
    try:
        mod = importlib.import_module(module_name)
        version = getattr(mod, '__version__', 'unknown')
        print(f"  ✓ {display_name:.<30} {version}")
        return True
    except ImportError as e:
        print(f"  ✗ {display_name:.<30} NOT INSTALLED")
        print(f"    Error: {e}")
        return False

def main():
    print("=" * 60)
    print("VideoDisco Phase 0 Setup Verification")
    print("=" * 60)
    print()
    
    print("🔍 Checking required Python packages:\n")
    
    all_passed = True
    
    # Core ML libraries
    print("ML & ONNX Runtime:")
    all_passed &= test_import('onnxruntime', 'onnxruntime-silicon')
    all_passed &= test_import('ultralytics', 'YOLOv8/ultralytics')
    all_passed &= test_import('torch', 'PyTorch')
    all_passed &= test_import('numpy', 'NumPy')
    
    print()
    print("MLOps & Monitoring:")
    all_passed &= test_import('mlflow', 'MLflow')
    
    print()
    print("Cloud Service:")
    all_passed &= test_import('fastapi', 'FastAPI')
    all_passed &= test_import('uvicorn', 'Uvicorn')
    
    print()
    print("Image & Computer Vision:")
    all_passed &= test_import('cv2', 'OpenCV')
    all_passed &= test_import('PIL', 'Pillow')
    
    print()
    print("Development & Testing:")
    all_passed &= test_import('pytest', 'Pytest')
    
    print()
    print("=" * 60)
    
    if all_passed:
        print("✅ Phase 0 Setup VERIFIED - All requirements installed!")
        print()
        print("Next steps:")
        print("  1. Run: python tests/test_onnx.py")
        print("  2. Test ONNX tensor inference")
        print("  3. Proceed to Phase 1: Edge Inference")
        return 0
    else:
        print("❌ Phase 0 Setup INCOMPLETE - Some packages missing")
        print()
        print("Install missing packages with:")
        print("  source venv/bin/activate")
        print("  pip install -r requirements.txt")
        return 1

if __name__ == "__main__":
    sys.exit(main())
