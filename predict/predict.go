package predict

import (
	"fmt"
	_ "image/png"
	"os"

	"github.com/metno/go-tf-cnn-predict/stddev"
	tf "github.com/wamuir/graft/tensorflow"
	"github.com/wamuir/graft/tensorflow/op"
)

func init() {
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "1")
	os.Setenv("KMP_AFFINITY", "noverbose")
}

func argmax(arr []float32) int {
	var big float32 = -1.0
	var idx int = -1

	for i := 0; i < len(arr); i++ {
		if arr[i] >= big {
			big = arr[i]
			idx = i
		}
	}
	return idx
}

type Predictor struct {
	Model     *tf.SavedModel
	InputName string
	ImgSize   int32
}

// declare method
func NewPredictor(modeldir string, inputName string, imgSize int32) (Predictor, error) {
	predictor := Predictor{}
	model, err := tf.LoadSavedModel(modeldir, []string{"serve"}, nil)
	if err != nil {
		fmt.Printf("Error loading saved model: %s\n", err.Error())
		return Predictor{}, err
	}
	predictor.Model = model
	predictor.InputName = inputName
	predictor.ImgSize = imgSize
	return predictor, err
}

// Predict - Return index with highest probability
func (p Predictor) Predict(imagefile string) (int, error) {
	float32arr, err := p.PredictionsArr(imagefile)
	if err != nil {
		return -1, err
	}
	cc := argmax(float32arr)
	return cc, err

}

// Predict - return class prediction and standard deviation
func (p Predictor) PredictWithVariance(imagefile string) (classpred int, deviation float64, err error) {
	float32arr, err := p.PredictionsArr(imagefile)
	if err != nil {
		return -1, -1, err
	}
	//fmt.Printf("%+v\n", float32arr)
	deviation = stddev.CalcVariance(float32arr)
	if len(float32arr) == 1 { // Binary classification
		if float32arr[0] <= 0.5 {
			classpred = 0
		} else {
			classpred = 1
		}
		deviation = float64(float32arr[0])
	} else { // multclass classification . (Or multilabel, handle eventually)
		classpred = argmax(float32arr)
	}

	return classpred, deviation, err

}

// Predict - return class prediction and standard deviation
func (p Predictor) PredictWithDeviation(imagefile string) (classpred int, deviation float64, err error) {
	float32arr, err := p.PredictionsArr(imagefile)
	if err != nil {
		return -1, -1, err
	}
	//fmt.Printf("%+v\n", float32arr)
	deviation = stddev.CalcDeviation(float32arr)
	if len(float32arr) == 1 { // Binary classification
		if float32arr[0] <= 0.5 {
			classpred = 0
		} else {
			classpred = 1
		}
		deviation = float64(float32arr[0])
	} else { // multclass classification . (Or multilabel, handle eventually)
		classpred = argmax(float32arr)
	}

	return classpred, deviation, err

}

// PredictionsArr Return array of probabilities
func (p Predictor) PredictionsArr(imagefile string) ([]float32, error) {

	//inputs := strings.Split(p.Model.Signatures["serving_default"].Inputs["densenet201_input"].Name, ":")
	//fmt.Printf("N: %+v\n", inputs[0])
	tensor, err := p.makeTensorFromImage(imagefile)
	if err != nil {
		return []float32{}, err

	}

	result, runErr := p.Model.Session.Run(
		map[tf.Output]*tf.Tensor{
			p.Model.Graph.Operation(p.InputName).Output(0): tensor,
		},
		[]tf.Output{
			p.Model.Graph.Operation("StatefulPartitionedCall").Output(0),
		},
		nil,
	)

	if runErr != nil {
		return []float32{}, runErr

	}

	arr32 := result[0].Value().([][]float32)[0]

	return arr32, nil
}

// Predict - return class prediction and standard deviation
func (p Predictor) PredictWithDeviationFromByteBufr(bytes []byte) (classpred int, deviation float64, err error) {
	float32arr, err := p.PredictionsArrFromByteBufr(bytes)
	if err != nil {
		return -1, -1, err
	}
	deviation = stddev.CalcDeviation(float32arr)
	if len(float32arr) == 1 { // Binary classification
		if float32arr[0] <= 0.5 {
			classpred = 0
		} else {
			classpred = 1
		}
		deviation = float64(float32arr[0])
	} else { // multclass classification . (Or multilabel, handle eventually)
		classpred = argmax(float32arr)
	}

	return classpred, deviation, err

}

// PredictionsArr Return array of probabilities
func (p Predictor) PredictionsArrFromByteBufr(bytes []byte) ([]float32, error) {

	tensor, err := p.makeTensorFromByteBufr(bytes)
	if err != nil {
		return []float32{}, err

	}

	result, runErr := p.Model.Session.Run(
		map[tf.Output]*tf.Tensor{
			p.Model.Graph.Operation(p.InputName).Output(0): tensor,
		},
		[]tf.Output{
			p.Model.Graph.Operation("StatefulPartitionedCall").Output(0),
		},
		nil,
	)

	if runErr != nil {
		return []float32{}, runErr

	}
	//fmt.Printf("RESULT: %+v\n", result[0].Value())
	arr32 := result[0].Value().([][]float32)[0]

	//fmt.Printf("Result: %v\n", result[0].Value().([][]float32))
	return arr32, nil
}

// Convert the image in filename to a Tensor suitable as input to the cc-classifier model.
func (p Predictor) makeTensorFromImage(filename string) (*tf.Tensor, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return p.makeTensorFromByteBufr(bytes)
}

// Convert the image in filename to a Tensor suitable as input to the cc-classifier model.
func (p Predictor) makeTensorFromByteBufr(bytes []byte) (*tf.Tensor, error) {

	// DecodeJpeg uses a scalar String-valued tensor as input.
	tensor, err := tf.NewTensor(string(bytes))
	if err != nil {
		return nil, err
	}
	// Construct a graph to normalize the image
	graph, input, output, err := p.constructGraphToNormalizeImage()
	if err != nil {
		return nil, err
	}

	// Execute that graph to normalize this one image
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}

	return normalized[0], nil

}

// The inception model takes as input the image described by a Tensor in a very
// specific normalized format (a particular image size, shape of the input tensor,
// normalized pixel values etc.).
//
// This function constructs a graph of TensorFlow operations which takes as
// input a JPEG-encoded string and returns a tensor suitable as input to the
// inception model.
func (p Predictor) constructGraphToNormalizeImage() (graph *tf.Graph, input, output tf.Output, err error) {

	// - The model was trained after with images scaled to 128x128 pixels.
	// - The colors, represented as R, G, B in 1-byte each were converted to
	//   float using (value - Mean)/Scale.
	const (
		Mean  = float32(1.0)
		Scale = float32(255)
	)
	// - input is a String-Tensor, where the string the JPEG-encoded image.
	// - The inception model takes a 4D tensor of shape
	//   [BatchSize, Height, Width, Colors=3], where each pixel is
	//   represented as a triplet of floats
	// - Apply normalization on each pixel and use ExpandDims to make
	//   this single image be a "batch" of size 1 for ResizeBilinear.
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	output = op.Div(s,
		op.Sub(s,
			op.ResizeBilinear(s,
				op.ExpandDims(s,
					op.Cast(s, decode, tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{p.ImgSize, p.ImgSize})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()

	return graph, input, output, err
}
