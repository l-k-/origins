from origins.exceptions import DoesNotExist, ValidationError
from . import neo4j, utils, deps
from .edges import Edge
from .components import Component


class Relationship(Edge):
    model_name = 'origins:Relationship'

    start_model = Component

    end_model = Component

    get_by_id_statement = '''

        MATCH (:`origins:Resource` {`origins:uuid`: { resource }})-[:manages]->(n:`origins:Relationship` {`origins:id`: { id }})
        RETURN n

    '''  # noqa

    # Gets the component's resource. This does not use the physical edge since
    # the component being looked up may be invalid
    get_resource_statement = '''

        MATCH (:`origins:Relationship` {`origins:uuid`: { uuid }})<-[:`origins:end`]-(:`origins:Edge` {`origins:type`: 'manages'})-[:`origins:start`]->(n:`origins:Resource`)
        RETURN n

    '''  # noqa

    def __init__(self, *args, **kwargs):
        self.resource = kwargs.pop('resource', None)
        super(Relationship, self).__init__(*args, **kwargs)

    @classmethod
    def get_by_id_query(cls, id, resource):
        statement = utils.prep(cls.get_by_id_statement,
                               model=cls.model_name)

        return {
            'statement': statement,
            'parameters': {
                'id': id,
                'resource': resource.uuid,
            }
        }

    @classmethod
    def _validate_unique(cls, node, tx):
        query = {
            'statement': cls.get_by_id_statement,
            'parameters': {
                'id': node.id,
                'resource': node.resource.uuid,
            }
        }

        result = tx.send(query)

        if result:
            raise ValidationError('{} already exists by id'
                                  .format(node.model_type))

    @classmethod
    def _add(cls, node, tx):
        prov = Edge._add(node, tx)

        # Define managing relationship
        Edge.add(start=node.resource,
                 end=node,
                 type='manages',
                 direction='bidirected',
                 dependence='inverse',
                 tx=tx)

        # Defining inclusion to resource
        Edge.add(start=node.resource,
                 end=node,
                 type='includes',
                 direction='bidirected',
                 dependence='none',
                 tx=tx)

        return prov

    @classmethod
    def _set(cls, old, new, tx):
        # Invalidate old version
        cls._invalidate(old, tx=tx)

        # Add new version
        prov = Edge._add(new, tx=tx)

        # Trigger the change dependency. This must occur after the new
        # node has been added so it is visible in the graph.
        deps.trigger_change(old, new, tx=tx)

        return prov

    @classmethod
    def resource(cls, uuid, tx=neo4j.tx):
        from .resources import Resource

        query = {
            'statement': cls.get_resource_statement,
            'parameters': {
                'uuid': uuid,
            }
        }

        result = tx.send(query)

        if not result:
            raise DoesNotExist('{} does not exist'.format(cls.model_type))

        return Resource.parse(*result[0])