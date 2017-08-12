/// <reference path="../References.d.ts"/>
import * as AgentTyeps from '../types/AgentTypes';

export function formatContinent(agent: AgentTyeps.Agent): string {
	return agent.continent && agent.continent_code ?
		agent.continent + ' (' + agent.continent_code + ')' :
		agent.continent || agent.continent_code || 'Unknown';
}

export function formatLocation(agent: AgentTyeps.Agent): string {
	return (agent.city ? agent.city + ', ' : '') +
		(agent.region || 'Unknown') +
		(agent.region_code ? ' (' + agent.region_code + ')' : '');
}

export function formatCountry(agent: AgentTyeps.Agent): string {
	return (agent.country || 'Unknown') +
		(agent.country_code ? ' (' + agent.country_code + ')' : '');
}

export function formatCoordinates(agent: AgentTyeps.Agent): string {
	return agent.latitude && agent.longitude ?
		agent.latitude + ', ' + agent.longitude : 'Unknown';
}
