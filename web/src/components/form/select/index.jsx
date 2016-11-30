import React from 'react'
import _ from 'lodash'
import { List } from 'immutable'

export default class FormSelect extends React.Component {
  constructor(props) {
    super(props)
  }

  onSelect(e) {
    const selected = this.props.options.findKey(option =>
      (option.get('value') === e.target.value)
    )

    this.props.onSelect(selected)
  }

  render() {
    const valueOfSelected = this.props.options.getIn([
      this.props.selected,
      'value',
    ])

    const noneOption = (!this.props.required) ? (
      <option value="">
        {this.props.defaultOption}
      </option>
    ) : null

    return (
      <div className="form-select">
        <select
          className="gray"
          onChange={this.onSelect.bind(this)}
          value={valueOfSelected}
        >
          {noneOption}

          {this.props.options.map((option, i) =>
            <option
              key={i}
              value={option.get('value')}
            >
              {option.get('label')}
            </option>
          )}
        </select>
      </div>
    )
  }
}

FormSelect.propTypes = {
  name: React.PropTypes.string.isRequired,
  defaultOption: React.PropTypes.string,
  options: React.PropTypes.instanceOf(List).isRequired,
  selected: React.PropTypes.number,
  required: React.PropTypes.bool,
  onSelect: React.PropTypes.func.isRequired,
}

FormSelect.defaultProps = {
  defaultOption: 'None',
  required: false,
}
