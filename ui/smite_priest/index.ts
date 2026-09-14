import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import { Sim } from '../core/sim.js';
import { TypedEvent } from '../core/typed_event.js';
import { SmitePriestSimUI } from './sim.js';

const sim = new Sim();
const player = new Player<Spec.SpecSmitePriest>(Spec.SpecSmitePriest, sim);
sim.raid.setPlayer(TypedEvent.nextEventID(), 0, player);

new SmitePriestSimUI(document.body, player);
