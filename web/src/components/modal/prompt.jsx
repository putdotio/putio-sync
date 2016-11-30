import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'

import Button from '../button'
import LinkButton from '../button/link'

import {
  Modal,
  ModalContent,
  ModalFooter,
  ModalFooterAction,
} from '../modal'

import {
  Form,
  Row,
  RowTitle,
  Input,
} from '../../components/form'

export default class Prompt extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  static Show(options) {
    return new Modal({
      name: 'prompt-modal',
      small: true,
      content: function() {
        return (
          <Prompt
            label={options.label}
            password={options.password}
            onCancel={(reason) => {
              this.Destroy(reason)
            }}
            onOkay={text => {
              this.Destroy(null, text)
            }}
          />
        )
      }
    }).Show()
  }

  onChange(value) {
    this.text = value
  }

  render() {
    return (
      <div>
        <Form>
          <Row>
            <RowTitle title={this.props.label} />
            <Input
              name="prompt-input"
              onChange={this.onChange.bind(this)}
              autofocus={true}
              password={this.props.password}
              onKeyUp={e => {
                if (e.which === 13) {
                  this.props.onOkay(this.text)
                }
              }}
            />
          </Row>
        </Form>

        <ModalFooter>
          <ModalFooterAction>
            <LinkButton
              onClick={() => this.props.onCancel(Modal.REASONS.CANCEL_BY_FOOTER) }
              label="Cancel"
            />
          </ModalFooterAction>

          <ModalFooterAction>
            <Button
              label="Okay"
              scope="btn-success"
              onClick={() => this.props.onOkay(this.text) }
            />
          </ModalFooterAction>
        </ModalFooter>
      </div>
    )
  }
}

Prompt.propTypes = {
  label: React.PropTypes.string.isRequired,
  password: React.PropTypes.bool.isRequired,
  onCancel: React.PropTypes.func.isRequired,
  onOkay: React.PropTypes.func.isRequired,
}

Prompt.defaultProps = {
  password: false,
}
