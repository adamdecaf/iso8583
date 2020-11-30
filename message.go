package iso8583

import (
	"fmt"

	"github.com/moov-io/iso8583/fields"
	"github.com/moov-io/iso8583/spec"
	"github.com/moov-io/iso8583/utils"
)

type Message struct {
	Fields map[int]fields.Field
	spec   *spec.MessageSpec

	// let's keep it 8 bytes for now
	bitmap *utils.Bitmap
}

func NewMessage(spec *spec.MessageSpec) *Message {
	return &Message{
		Fields: map[int]fields.Field{},
		spec:   spec,
	}
}

func (m *Message) Set(id int, field fields.Field) {
	m.Fields[id] = field
}

func (m *Message) Field(id int, val string) {
	m.Fields[id] = fields.NewField(id, []byte(val))
}

func (m *Message) BinaryField(id int, val []byte) {
	m.Fields[id] = fields.NewField(id, val)
}

func (m *Message) GetString(id int) string {
	// check index
	return m.Fields[id].String()
}

func (m *Message) GetBytes(id int) []byte {
	// check index
	return m.Fields[id].Bytes()
}

func (m *Message) Pack() ([]byte, error) {
	packed := []byte{}

	m.bitmap = utils.NewBitmap(128)

	// we need to test if max id will still be 0 or 1
	maxId := 0
	for id, _ := range m.Fields {
		if id > maxId {
			maxId = id
		}
		// regular fields start from index 2
		if id < 2 {
			continue
		}
		m.bitmap.Set(id)
	}

	// depending on the max Id generate bitmap:
	// 64, 128 or 192 bits

	// pack MTI
	packedMTI, err := m.spec.Fields[0].Pack(m.Fields[0].Bytes())
	if err != nil {
		return nil, err
	}
	packed = append(packed, packedMTI...)

	packedBitmap, err := m.spec.Fields[1].Pack(m.bitmap.Bytes())

	if err != nil {
		return nil, err
	}
	packed = append(packed, packedBitmap...)

	for i := 2; i <= maxId; i++ {
		// check if i exist
		if field, ok := m.Fields[i]; ok {
			def, ok := m.spec.Fields[i]
			if !ok {
				return nil, fmt.Errorf("Failed to pack field: %d. No definition found", i)
			}
			packedField, err := def.Pack(field.Bytes())
			if err != nil {
				return nil, err
			}
			packed = append(packed, packedField...)
		}
	}

	// m.packer.Pack(m, spec)
	// go through each spec field
	// add packed MTI
	// add packed bitmap (find bitmap definition - field N1)
	// pack each field starting from 2
	// return result

	// for id, fieldPacker := range m.spec.Fields {
	// }
	// packer := &packer.MessagePacker{}

	// strange argument passing :)
	// packer.Pack(m, m.Spec)

	return packed, nil
}

func (m *Message) Unpack(src []byte) error {
	return nil
}
