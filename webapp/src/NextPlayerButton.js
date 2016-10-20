import React, {Component} from 'react'
import {FormattedMessage} from 'react-intl';

class NextPlayerButton extends Component {

    nextPlayer() {
        fetch('/api/games/' + this.props.gameId + '/hold', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({})
        })
            .then((response) => response.json())
            //.then((json) => browserHistory.push("/games/" + json.id))
            .catch((error) => console.log(error))
    }

    render() {
        const label = ( this.props.game.Ongoing == 3 ? <FormattedMessage id='nextPlayerBtn.nextBtn.label' defaultMessage='Next'/> : <FormattedMessage id='nextPlayerBtn.holdBtn.label' defaultMessage='Hold'/> )
        return (
            <div>
                <a className="waves-effect waves-light btn light-blue" onClick={() => this.nextPlayer() }>{label}</a>
            </div>
        )
    }
}

export default NextPlayerButton
