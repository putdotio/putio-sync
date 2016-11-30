import React from 'react'
import Button from './index'

export default class LinkButton extends React.Component {
  render() {
    return (
      <Button
        {...this.props}
        link={true}
      />
    )
  }
}
