from flask import Flask, request, jsonify
import pytesseract
from PIL import Image
import io

app = Flask(__name__)

@app.route('/ocr', methods=['POST'])
def ocr_image():
    if 'image' not in request.files:
        return jsonify({'error': 'No image uploaded'}), 400

    image_file = request.files['image']
    try:
        image = Image.open(image_file.stream)
    except Exception as e:
        return jsonify({'error': f'Invalid image: {str(e)}'}), 400

    # OCR with English and Russian
    text = pytesseract.image_to_string(image, lang='eng+rus')

    return jsonify({'text': text})

if __name__ == '__main__':
    app.run(debug=True)
