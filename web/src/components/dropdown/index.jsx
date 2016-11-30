import React from 'react'
import PureRenderMixin from 'react-addons-pure-render-mixin'
import { List } from 'immutable'
import _ from 'lodash'
import Button from '../button'
import { Select } from '../../components/form'

export default class Dropdown extends React.Component {
  constructor() {
    super()
    this.shouldComponentUpdate = PureRenderMixin.shouldComponentUpdate.bind(this)
  }

  render() {
    const selected = this.props.options.get(this.props.selected)

    if (this.props.isMobile) {
      return (
        <Select
          name="mobile-dropdown"
          options={this.props.options}
          selected={this.props.selected}
          defaultOption={this.props.label}
          onSelect={index => {
            this.props.onChange(this.props.options.get(index), index)
          }}
        />
      )
    }

    return (
      <div className="dropdown dropdown-gray">
        <Button
          onClick={() => {
            if (this.props.onClick) {
              this.props.onClick(this.props.options.get(0), 0)
            }
          }}
          icon={this.props.icon}
          iconRight={this.props.iconRight}
          label={(selected) ? selected.get('label') : this.props.label}
          scope={this.props.scope || "btn-default"}
          fixed={this.props.fixed || false}
        />

        <div className="dropdown-content">
          {this.props.options.map((option, i) => {
            const iconLeft = (option.get('icon')) ? (
              <i className={option.get('icon')}></i>
            ) : null

            const iconRight = (option.get('iconRight')) ? (
              <i className={option.get('iconRight')}></i>
            ) : null

            const withIcon = (iconLeft || iconRight)

            const className = _.compact([
              'dropdown-option',
              (withIcon) ? 'with-icon' : null,
              this.props.scope || 'btn-default',
            ]).join(' ')

            return (
              <div
                className={className}
                key={i}
                onClick={() => this.props.onChange(option, i)}
              >
                <a>
                  <span>
                    {iconLeft}
                    {option.get('label')}
                    {iconRight}
                  </span>
                </a>
              </div>
            )
          })}
        </div>
      </div>
    )
  }
}

Dropdown.propTypes = {
  options: React.PropTypes.instanceOf(List),
  selected: React.PropTypes.number,
  onChange: React.PropTypes.func.isRequired,
  onClick: React.PropTypes.func,
  icon: React.PropTypes.string,
  isMobile: React.PropTypes.bool,
  scope: React.PropTypes.string,
  fixed: React.PropTypes.bool,
}

Dropdown.defaultProps = {
  isMobile: false,
}
