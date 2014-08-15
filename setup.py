from setuptools import setup, find_packages


kwargs = {
    'name': 'origins',
    'version': __import__('origins').get_version(),
    'description': 'Origins is a free and open source service for building '
                   'information dependency graphs across your data, '
                   'systems, and operations.',
    'url': 'https://github.com/cbmi/origins/',

    'packages': find_packages(exclude=['tests']),

    'author': 'Byron Ruth',
    'author_email': 'b@devel.io',

    'license': 'BSD',
    'keywords': 'graph dependency provenance information service REST',
    'classifiers': [
        'Development Status :: 3 - Alpha',
        'License :: OSI Approved :: BSD License',
        'Intended Audience :: Developers',
        'Intended Audience :: Science/Research',
    ],

    'install_requires': [
        'requests>=2.2.0,<2.3',
        'flask>=0.10,<0.11',
        'flask-restful>=0.2,<0.3',
        'flask-cors>=1.7,<1.8',
    ],

    'test_suite': 'test_suite',
}

setup(**kwargs)
